package log

import (
	"os"
	"sync"
)

const version = "0.0.0"

// init routines to manage log processing.
func init() {
	l.tags = newTags()
	l.conf = defaultConfig()
	l.errFile, _ = os.Create("/tmp/blitz.log")

	go local() // start local service
	go mux()   // start multiplexer
}

type logging struct {
	conf    *config
	wg      sync.WaitGroup
	errFile *os.File
	tags    *tags
}

var l logging
