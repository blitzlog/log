package log

import (
	"github.com/blitzlog/proto/log"
)

// declare channels for piping log messages
var (
	logChannel   = make(chan *log.Log, 1000) // input log channel
	edgeChannel  = make(chan *log.Log, 1000) // channel to push logs to edge
	localChannel = make(chan *log.Log, 1000) // channel to publish logs locally
	flushChannel = make(chan bool, 1)        // channel to flush logs
)

// mux send logs on multiple channels based on config.
func mux() {
	for lg := range logChannel {
		// log local if API key not set, or there is an API error
		logLocal := l.conf.apiKey == "" || l.conf.apiError
		switch {
		case logLocal:
			localChannel <- lg
		default:
			edgeChannel <- lg
		}
	}
}
