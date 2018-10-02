package log

import (
	"fmt"

	"github.com/blitzlog/proto/log"
)

type config struct {
	logLevel     log.Level // current log type
	logVerbosity int32     // current log level
	logJson      bool      // log as json
	apiKey       string    // API Key
	apiError     bool      // API Key is incorrect
}

func defaultConfig() *config {
	return &config{}
}

func SetAPIKey(key string) {
	l.conf.apiKey = key
	redirect()
	sender()
}

func JSON() {
	l.conf.logJson = true
}

func String(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("%s", v)
	}
}

// Define log level constants, used for setting log level by user.
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

func SetLevel(level string) {
	l.conf.logLevel = log.Level(log.Level_value[level])
}

func GetLevel() string {
	return l.conf.logLevel.String()
}

// Verbosity records the verbosity of a log.
type Verbosity int32

var defaultVerbosity = V(0)

// V creates new verbosity.
func V(verbosity int32) *Verbosity {
	v := Verbosity(verbosity)
	return &v
}

// V updates verbosity
func (vtags *VTags) V(verbosity int32) *VTags {
	vtags.v = V(verbosity)
	return vtags
}

// verbose checks if it is verbose as current config.
func (v *Verbosity) verbose() bool {
	return int32(*v) <= l.conf.logVerbosity
}

func SetVerbosity(v int32) {
	l.conf.logVerbosity = v
}

func GetVerbosity() int32 {
	return l.conf.logVerbosity
}
