package log

import (
	"github.com/blitzlog/proto/log"
)

// mux log to local and/or edge.
func mux(lg *log.Log) {

	// log local if
	// - API key not set
	// - config set to log local
	// - error sending log to edge
	if l.conf.apiKey == "" || l.conf.logLocal || l.conf.apiError {
		logLocal(lg)
	}

	// log edge if api key is set and no errors sending to edge.
	if l.conf.apiKey != "" && !l.conf.apiError {
		l.wg.Add(1)
		l.edgeChannel <- lg
	}
}
