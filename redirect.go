package log

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/blitzlog/proto/log"
)

// redirect stdout to log channel.
func redirect() {

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	reader := bufio.NewReader(r)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				l.errFile.WriteString(err.Error())
			}
			l.wg.Add(1)
			l.logChannel <- &log.Log{
				Timestamp: time.Now().UTC().UnixNano() / 1e6,
				Raw:       strings.TrimSpace(line),
			}
		}
	}()
}
