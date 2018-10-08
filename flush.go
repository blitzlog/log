package log

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/blitzlog/proto/log"
)

// Flush all logs sent so far.
func Flush() {

	// flush stdout
	time.Sleep(time.Millisecond)
	l.stdout.Sync()

	// if we are emitting logs, then get stack trace
	logLocal := l.conf.apiKey == "" || l.conf.apiError
	if !logLocal {
		r := recover()
		if r != nil {
			stack := debug.Stack()
			l.wg.Add(1)
			logChannel <- &log.Log{
				Timestamp: time.Now().UTC().UnixNano() / 1e6,
				Raw:       fmt.Sprintf("%s\n%s", r, stack),
			}
		}
	}

	// flush logs
	time.Sleep(time.Millisecond)
	flushChannel <- true

	// wait to process all logs
	l.wg.Wait()
}
