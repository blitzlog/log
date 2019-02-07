package log

import (
	"os"
	"sync"

	"github.com/blitzlog/proto/log"
)

const version = "0.0.0"

// init routines to manage log processing.
func init() {
	l.tags = newTags()
	l.conf = defaultConfig()
	l.errFile, _ = os.Create("/tmp/blitz.log")
	l.stdout = os.Stdout

	// init channels
	l.logChannel = make(chan *log.Log, 1000)
	l.edgeChannel = make(chan *log.Log, 1000)
	l.localChannel = make(chan *log.Log, 1000)
	l.flushChannel = make(chan bool, 1)

	go local() // start local service
	go mux()   // start multiplexer

	// TODO: enable configurable stdout redirect
	//redirect() // redirect logs from stdout
}

type logging struct {
	conf         *config
	wg           sync.WaitGroup
	stdout       *os.File
	errFile      *os.File
	tags         *tags
	logChannel   chan *log.Log // input log channel
	edgeChannel  chan *log.Log // channel to push logs to edge
	localChannel chan *log.Log // channel to publish logs locally
	flushChannel chan bool     // channel to flush logs
}

var l logging
