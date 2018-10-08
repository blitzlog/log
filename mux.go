package log

// mux send logs on multiple channels based on config.
func mux() {
	for lg := range l.logChannel {
		// log local if API key not set, or there is an API error
		logLocal := l.conf.apiKey == "" || l.conf.apiError
		switch {
		case logLocal:
			l.localChannel <- lg
		default:
			l.edgeChannel <- lg
		}
	}
}
