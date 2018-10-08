package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/blitzlog/proto/log"
)

const ErrorKey = "error"

// With adds tags to log.
func (v *Verbosity) With(tags Tags) *VTags {
	return &VTags{v, tags}
}

func (v *Verbosity) Tag(key string, val interface{}) *VTags {
	tags := Tags{key: val}
	return &VTags{v, tags}
}

func (v *Verbosity) D(format string, args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_debug, nil, format, args)
	}
}

func (v *Verbosity) I(format string, args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_info, nil, format, args)
	}
}

func (v *Verbosity) W(format string, args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_warn, nil, format, args)
	}
}

func (v *Verbosity) E(format string, args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_error, nil, format, args)
	}
}

func (v *Verbosity) F(format string, args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_fatal, nil, format, args)
	}
}

func (v *Verbosity) Debug(args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_debug, nil, fmt.Sprint(args...), nil)
	}
}

func (v *Verbosity) Info(args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_info, nil, fmt.Sprint(args...), nil)
	}
}

func (v *Verbosity) Warn(args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_warn, nil, fmt.Sprint(args...), nil)
	}
}

func (v *Verbosity) Error(args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_error, nil, fmt.Sprint(args...), nil)
	}
}

func (v *Verbosity) Fatal(args ...interface{}) {
	if v.verbose() {
		pushLog(v, log.Level_fatal, nil, fmt.Sprint(args...), nil)
	}
}

type Tags map[string]interface{}

func (tags Tags) stringTags() map[string]string {
	strTags := make(map[string]string)
	for k, v := range tags {
		strTags[k] = String(v)
	}
	return strTags
}

type VTags struct {
	v    *Verbosity
	tags Tags
}

// With adds tags to log.
func With(tags Tags) *VTags {
	return &VTags{defaultVerbosity, tags}
}

func Tag(k string, v interface{}) *VTags {
	return &VTags{defaultVerbosity, Tags{k: v}}
}

// With add tags to vtag.
func (vtags *VTags) With(tags Tags) *VTags {
	for k, v := range tags {
		vtags.tags[k] = v
	}
	return vtags
}

func (vtags *VTags) Tag(k string, v interface{}) *VTags {
	vtags.tags[k] = v
	return vtags
}

func WithError(err error) *VTags {
	tags := map[string]interface{}{ErrorKey: err}
	return &VTags{defaultVerbosity, tags}
}

func (vtags *VTags) WithError(err error) *VTags {
	vtags.tags[ErrorKey] = err
	return vtags
}

func (v *Verbosity) WithError(err error) *VTags {
	tags := map[string]interface{}{ErrorKey: err}
	return &VTags{v, tags}
}

func (vtags *VTags) D(format string, args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_debug, vtags.tags,
			format, args)
	}
}

func (vtags *VTags) I(format string, args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_info, vtags.tags,
			format, args)
	}
}

func (vtags *VTags) W(format string, args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_warn, vtags.tags,
			format, args)
	}
}

func (vtags *VTags) E(format string, args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_error, vtags.tags,
			format, args)
	}
}

func (vtags *VTags) F(format string, args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_fatal, vtags.tags,
			format, args)
	}
}

func (vtags *VTags) Debug(args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_debug, vtags.tags,
			fmt.Sprint(args...), nil)
	}
}

func (vtags *VTags) Info(args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_info, vtags.tags,
			fmt.Sprint(args...), nil)
	}
}

func (vtags *VTags) Warn(args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_warn, vtags.tags,
			fmt.Sprint(args...), nil)
	}
}

func (vtags *VTags) Error(args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_error, vtags.tags,
			fmt.Sprint(args...), nil)
	}
}

func (vtags *VTags) Fatal(args ...interface{}) {
	if vtags.v.verbose() {
		pushLog(defaultVerbosity, log.Level_fatal, vtags.tags,
			fmt.Sprint(args...), nil)
	}
}

func D(format string, args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_debug, nil,
		format, args)
}

func I(format string, args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_info, nil,
		format, args)
}

func W(format string, args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_warn, nil,
		format, args)
}

func E(format string, args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_error, nil,
		format, args)
}

func F(format string, args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_fatal, nil,
		format, args)
}

func Debug(args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_debug, nil,
		fmt.Sprint(args...), nil)
}

func Info(args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_info, nil,
		fmt.Sprint(args...), nil)
}

func Warn(args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_warn, nil,
		fmt.Sprint(args...), nil)
}

func Error(args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_error, nil,
		fmt.Sprint(args...), nil)
}

func Fatal(args ...interface{}) {
	pushLog(defaultVerbosity, log.Level_fatal, nil,
		fmt.Sprint(args...), nil)
}

// pushLog creates a Log object and pushes it over the encodeChannel.
func pushLog(verbosity *Verbosity, level log.Level, tags Tags,
	format string, args []interface{}) {

	// check if this log type is to be logged
	if level < l.conf.logLevel {
		return
	}

	// increment wait group
	l.wg.Add(1)

	// get location info for the log
	file, function, line := fileLine(3)

	l.logChannel <- &log.Log{
		File:      file,
		Line:      int32(line),
		Function:  function,
		Timestamp: time.Now().UTC().UnixNano() / 1e6,
		Level:     level,
		Verbosity: int32(*verbosity),
		Msg:       fmt.Sprintf(format, args...),
		Tags:      tags.stringTags(),
	}
	if level == log.Level_fatal {
		Flush()
		panic(fmt.Sprintf(format, args...))
	}
}

// fileLine returns the file, function and line for calling function.
func fileLine(depth int) (string, string, int) {

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "???", "???", 1
	}

	// prune file name
	slash := strings.LastIndex(file, "/")
	if slash >= 0 {
		file = file[slash+1:]
	}

	// get function name
	fn := runtime.FuncForPC(pc).Name()
	slash = strings.LastIndex(fn, "/")
	if slash >= 0 {
		fn = fn[slash+1:]
	}
	slash = strings.LastIndex(fn, ".")
	if slash >= 0 {
		fn = fn[slash+1:]
	}

	return file, fn, line
}
