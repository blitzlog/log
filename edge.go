package log

// Send encoded logs to edge server.

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/blitzlog/errors"
	"github.com/blitzlog/proto/edge"
	"github.com/blitzlog/proto/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server address
const address = "test.blitzlog.com:8089"

// retryLimit at each step when sending logs to edge server.
const retryLimit = 4

// Tx trasmits messages to edge server.
// Exported to enable unit test of api server.
type Tx struct {
	token      string
	edgeClient edge.EdgeClient
	logClient  edge.Edge_PostLogsClient
	logMap     map[string]int32
	latency    int32
	errCount   int32
	retryCount int
}

func NewTx() *Tx {
	return &Tx{logMap: make(map[string]int32)}
}

// sender reads messages from senderChannel and sends to edge server.
func sender() {

	// create new transmitter
	tx := NewTx()

	// initialize transmitter
	var lgs []*log.Log

	// initialize pause counter
	var pause int

	// initialize flush ticker
	flushDuration := time.Second
	flushTick := time.NewTicker(flushDuration)

	// accumulate and send logs
	go func() {
		for {
			select {
			case <-flushTick.C:
				if pause == 0 {
					lgs, pause = tx.send(lgs)
					continue
				}
				pause--
			case lg := <-edgeChannel:
				lgs = append(lgs, lg)
			}
		}
	}()
}

// send logs to edge client, with exponential backtracking in case of failures.
func (tx *Tx) send(lgs []*log.Log) ([]*log.Log, int) {

	defer func() {
		l.errFile.Sync()
	}()

	var err error

	// create edge client if does not exist
	if tx.edgeClient == nil {
		tx.edgeClient, err = getEdgeClient()
		if err == nil {
			tx.retryCount = 0
		}
	}

	// handle edge client error
	if err != nil {
		l.errFile.WriteString(fmt.Sprintf("edge client error: %v\n", err))
		tx.errCount++
		tx.retryCount++
		return lgs, 2 ^ (tx.retryCount - 1)
	}

	// create token if empty
	if tx.token == "" {
		startMs := nowMs()
		tx.token, err = getToken(tx.edgeClient, l.conf.apiKey)
		tx.latency = int32(nowMs() - startMs)

		// clear retry count
		if err == nil {
			tx.retryCount = 0
		}
	}

	// handle get token error
	if err != nil {
		l.errFile.WriteString(fmt.Sprintf("token error: %v\n", err))

		// if at retry limit then backtrack
		if tx.retryCount == retryLimit {
			l.errFile.WriteString("backtracking to edge client\n")
			tx.edgeClient = nil
			tx.retryCount = 0
		}
		tx.errCount++
		tx.retryCount++
		return lgs, 2 ^ (tx.retryCount - 1)
	}

	// create log client
	if tx.logClient == nil {
		tx.logClient, err = tx.edgeClient.PostLogs(context.Background())
		if err == nil {
			tx.logMap = make(map[string]int32)
			tx.retryCount = 0
		}
	}

	// handle log client error
	if err != nil {
		l.errFile.WriteString(fmt.Sprintf("log client error: %v\n", err))

		// if at retry limit then backtrack
		if tx.retryCount == retryLimit {
			l.errFile.WriteString("backtracking to get token\n")
			tx.token = ""
			tx.retryCount = 0
		}
		tx.errCount++
		tx.retryCount++
		return lgs, 2 ^ (tx.retryCount - 1)
	}

	// aggregate logs
	logs := new(log.Logs)
	for _, lg := range lgs {
		logs = tx.Append(logs, lg)
	}

	// send logs
	tx.latency, err = sendLogs(tx.logClient, tx.token, logs, tx.latency, tx.errCount)

	// handle send log errors
	if err != nil {
		l.errFile.WriteString(fmt.Sprintf("error sending logs: %v\n", err))

		// if at retry limit then backtrack
		if tx.retryCount == retryLimit {
			l.errFile.WriteString("backtracking to get log client\n")
			tx.logClient = nil
			resetGlobalTags()
			tx.retryCount = 0
		}

		tx.errCount++
		tx.retryCount++
		return lgs, 2 ^ (tx.retryCount - 1)
	}

	// update error and retry count
	tx.errCount = 0
	tx.retryCount = 0

	return nil, 2 ^ 0
}

// Append log to encoded logs.
func (tx *Tx) Append(logs *log.Logs, lg *log.Log) *log.Logs {

	// if raw log then append to raws
	if lg.GetRaw() != "" {
		logs.Raws = append(logs.Raws,
			&log.LogRaw{
				Timestamp: lg.GetTimestamp(),
				Raw:       lg.GetRaw(),
			})
		return logs
	}

	// else encode and add
	logKey, logVal := splitLog(lg)
	lookupKey := getLookupKey(logKey)

	index, ok := tx.logMap[lookupKey]
	if !ok {
		index = int32(len(tx.logMap))
		tx.logMap[lookupKey] = index
		logs.Keys = append(logs.Keys, logKey)
	}
	logVal.Index = index
	logs.Vals = append(logs.Vals, logVal)

	return logs
}

// sendLogs to edge server, if any logs or tags received.
// Tracks latency and error count of messages to edge server. Each message
// includes latency for, and count of errors since last last successful
// message sent to edge server.
func sendLogs(logClient edge.Edge_PostLogsClient, token string,
	logs *log.Logs, latency, errCount int32) (int32, error) {

	// get global tags
	logs.InstTags = getGlobalTags()

	// return if nothing to send
	if len(logs.Vals) == 0 && len(logs.InstTags) == 0 && len(logs.Raws) == 0 {
		return latency, nil
	}

	// create edge metrics
	metrics := &edge.Metrics{
		Latency:          latency,
		ErrCount:         errCount,
		LogChannelSize:   int32(len(logChannel)),
		EdgeChannelSize:  int32(len(edgeChannel)),
		LocalChannelSize: int32(len(localChannel)),
	}

	// create post logs request
	req := &edge.PostLogsRequest{
		TokenId: token,
		Logs:    logs,
		Metrics: metrics,
	}

	startMs := nowMs()
	err := logClient.Send(req)
	if err != nil {
		return latency, errors.Wrap(err, "send error")
	}

	resp, err := logClient.Recv()
	if err != nil {
		return latency, errors.Wrap(err, "response error")
	}

	if resp.Code != 200 {
		return latency, errors.New("grpc response: %d", resp.Code)
	}

	// update log level and verbosity based on response
	if resp.GetLogLevel() != log.Level_none {
		SetLevel(resp.GetLogLevel().String())
	}
	// verbosity is encoded as +1, so we subtract 1 and apply
	if resp.GetLogVerbosity() != 0 {
		SetVerbosity(resp.GetLogVerbosity() - 1)
	}

	// calculate latency for sending logs to edge server
	latency = int32(nowMs() - startMs)

	// update wait group for each log sent
	for _ = range logs.Vals {
		l.wg.Done()
	}
	for _ = range logs.Raws {
		l.wg.Done()
	}

	return latency, nil
}

func splitLog(lg *log.Log) (*log.LogKey, *log.LogVal) {
	return &log.LogKey{
			File:      lg.File,
			Line:      lg.Line,
			Function:  lg.Function,
			Level:     lg.Level,
			Verbosity: lg.Verbosity,
			Msg:       lg.Msg,
		}, &log.LogVal{
			Timestamp: lg.Timestamp,
			LineTags:  lg.Tags,
		}
}

func getLookupKey(logKey *log.LogKey) string {
	return fmt.Sprintf("%s:%d:%s:%s",
		logKey.File, logKey.Line, logKey.Function, logKey.Msg)
}

// getCredentials uses hardcoded certificate to create TLS credentials
// that would be used to connect to edge server.
func getCredentials() (credentials.TransportCredentials, error) {
	b := []byte(serverCert)
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{RootCAs: cp}), nil
}

// getEdgeClient creates new edge client.
func getEdgeClient() (edge.EdgeClient, error) {

	// DEBUG: use debug connector for logging dialer errors.
	//conn, err := debugConn()

	creds, err := getCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "error getting credentials")
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "error dialing to server")
	}

	return edge.NewEdgeClient(conn), nil
}

// debugConn creats a grpc connection that logs dialer errors.
func debugConn() (*grpc.ClientConn, error) {

	creds, err := getCredentials()
	if err != nil {
		return nil, errors.Wrap(err, "error getting credentials")
	}

	ctx := context.TODO()

	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		conn, err := (&net.Dialer{Cancel: ctx.Done()}).Dial("tcp", address)
		if err != nil {
			return nil, errors.Wrap(err, "dial error")
		}
		if creds != nil {
			conn, _, err = creds.ClientHandshake(ctx, address, conn)
			if err != nil {
				return nil, errors.Wrap(err, "handshake error")
			}
		}
		return conn, nil
	}

	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithDialer(dialer),
		grpc.WithInsecure(), //dialer handles TLS
	)

	return grpc.Dial(address, opts...)
}

// getToken uses API key to get a token from edge server.
func getToken(c edge.EdgeClient, keyId string) (string, error) {

	authRequest := &edge.AuthRequest{
		Version: version,
		KeyId:   keyId,
	}

	authResponse, err := c.Authenticate(context.Background(), authRequest)
	if err != nil {
		return "", errors.Wrap(err, "error authenticating")
	}

	if authResponse.Code != 200 {
		return "", errors.New("unauthorized request")
	}

	return authResponse.GetTokenId(), nil
}

// nowMs returns current time in milliseonds.
func nowMs() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
