package log

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/blitzlog/proto/log"
)

// Flush all logs sent so far.
func Flush() {

	// sleep to flush all stdout logs
	time.Sleep(time.Millisecond)

	// if we are emitting logs, then get stack trace
	logLocal := l.conf.apiKey == "" || l.conf.apiError
	if !logLocal {
		r := recover()
		if r != nil {
			l.errFile.WriteString(fmt.Sprintf("recover: %v", r))
			stack := debug.Stack()
			l.wg.Add(1)
			logChannel <- &log.Log{
				Timestamp: time.Now().UTC().UnixNano() / 1e6,
				Raw:       string(stack),
			}
		}
	}

	// wait to process all logs
	l.wg.Wait()
}
