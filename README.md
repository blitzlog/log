# Blitzlog

Structured log publishing library.

```
package main

import (
	"github.com/blitzlog/log"
)

const data = "some data"

func main() {
	defer log.Flush()                                            // flush all logs before exit

	log.SetVerbosity(2)                                          // logs till this verbosity are printed

	log.I("an info log")                                         // simple info log

	log.D("add data: %s", data)                                  // log can be formatted

	log.With(log.Tags{"tag": "value"}).W("tagged warning log")   // add tags to log

	log.With(log.Tags{"t1": 1, "t2": data}).I("multiple tags")   // tags should implement Stringer interface

	log.V(2).I("log prints at verbosity 2 or more")              // set verbosity of the log

	log.With(log.Tags{"t": "v"}).V(1).I("all of it together")    // all together

	log.V(1).With(log.Tags{"t": "v"}).I("flip the order around") // flip the order around

	log.F("going to panic")                                      // fatal terminate program execution

	log.I("never gets to this log line")                         // log line never executed
}
```

* Print JSON logs, default is concise human readable format.
	* `log.JSON()`
* Change minimum log level to be published
	* `log.SetLevel(log.InfoLevel)`
	* Levels in increasing order `Debug`, `Info`, `Warn`, `Error`, and `Fatal`.
* Change maximum log verbosity to be published
	* `log.SetVerbosity(2)`

### verbosity?

Each log line has an associated verbosity - a positive interger, which is `0` by default. Also there is a `global log verbosity`, only logs with verbosity less than or equal to the global verbosity are published. Assigning higher verbosity to more detailed logs helps control log volume.

### why defer?

This log publishing library may be used to push logs to a log server (or stdout), `defer Flush()` enables cleanly pushing logs over network even in case of panic, while maintaining lightining fast speeds.

Logs may be published from a local device, private or public clouds (AWS, GCP, Azure, Digital Ocean...), and from various kinds of deployments including container based (ECS, K8, Mesos) ones.
