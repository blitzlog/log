package log

import (
	"sync"
)

// tags keps track of global tags.
type tags struct {
	mu    sync.Mutex
	reset bool
	all   map[string]string
	dirty map[string]string
}

// newTags initializes tags structure.
func newTags() *tags {
	return &tags{
		all:   make(map[string]string),
		dirty: make(map[string]string),
	}
}

// Global sets global tags, for all logs sent via this instance.
func Global(tags Tags) {
	l.tags.mu.Lock()
	defer l.tags.mu.Unlock()
	for k, v := range tags {
		l.tags.all[k] = String(v)
		l.tags.dirty[k] = String(v)
	}
}

// getGlobalTags returns new global tags.
func getGlobalTags() map[string]string {

	l.tags.mu.Lock()
	defer l.tags.mu.Unlock()

	var tags map[string]string

	switch {
	case l.tags.reset:
		tags = l.tags.all
	default:
		tags = l.tags.dirty
	}

	l.tags.reset = false

	if len(tags) == 0 {
		return nil
	}

	l.tags.dirty = make(map[string]string)
	return tags
}

// resetGlobalTags forces re-sending all global tags,
// used when log client connection breaks.
func resetGlobalTags() {
	l.tags.mu.Lock()
	defer l.tags.mu.Unlock()
	l.tags.reset = true
}
