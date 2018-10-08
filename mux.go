package log

// mux send logs on multiple channels based on config.
func mux() {
	for lg := range l.logChannel {

		// only local if API key not set, or there is an API error
		onlyLocal := l.conf.apiKey == "" || l.conf.apiError

		// only edge if not log local and log.Local() not specified
		onlyEdge := !onlyLocal && !l.conf.logLocal

		switch {
		case onlyLocal:
			l.localChannel <- lg
		case onlyEdge:
			l.edgeChannel <- lg
		default:
			l.wg.Add(1)
			l.localChannel <- lg
			l.edgeChannel <- lg
		}
	}
}
