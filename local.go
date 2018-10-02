package log

import (
	"fmt"
	"time"

	"github.com/blitzlog/proto/log"
)

func local() {
	for lg := range localChannel {

		switch {
		case l.conf.logJson:
			fmt.Println(jsonFormat(lg))
		default:
			fmt.Println(format(lg))
		}
		l.wg.Done()
	}
}

// format log as:
// TMMDD HH:MM:SS.sss file:line <msg> <k1=v1 k2=v2>
func format(lg *log.Log) string {

	var buf []byte

	switch lg.Level {
	case log.Level_debug:
		buf = []byte("D")
	case log.Level_info:
		buf = []byte("I")
	case log.Level_warn:
		buf = []byte("W")
	case log.Level_error:
		buf = []byte("E")
	case log.Level_fatal:
		buf = []byte("F")
	}

	ts := time.Unix(0, lg.GetTimestamp()*int64(time.Millisecond)).Format("0102 15:04:05.000")
	buf = append(buf, []byte(ts)...)
	buf = append(buf, []byte(" ")...)
	buf = append(buf, []byte(lg.GetFile())...)
	buf = append(buf, []byte(":")...)
	buf = append(buf, []byte(fmt.Sprintf("%d", lg.GetLine()))...)
	buf = append(buf, []byte(" ")...)
	buf = append(buf, []byte(lg.GetMsg())...)
	tags := lg.GetTags()
	for k, v := range tags {
		buf = append(buf, fmt.Sprintf(" %s=%s", k, v)...)
	}

	return string(buf)
}

// jsonFormat log.
func jsonFormat(lg *log.Log) string {

	buf := []byte("{")

	// TODO: use String()
	switch lg.Level {
	case log.Level_debug:
		buf = append(buf, []byte("\"type\":\"debug\"")...)
	case log.Level_info:
		buf = append(buf, []byte("\"type\":\"info\"")...)
	case log.Level_warn:
		buf = append(buf, []byte("\"type\":\"warn\"")...)
	case log.Level_error:
		buf = append(buf, []byte("\"type\":\"error\"")...)
	case log.Level_fatal:
		buf = append(buf, []byte("\"type\":\"fatal\"")...)
	}

	ts := time.Unix(0, lg.GetTimestamp()*int64(time.Millisecond)).Format("2006-01-02 15:04:05.000")
	tsStr := fmt.Sprintf(", \"timestamp\":\"%s\"", ts)
	buf = append(buf, []byte(tsStr)...)
	buf = append(buf, []byte(", \"file\":\"")...)
	buf = append(buf, []byte(lg.GetFile())...)
	buf = append(buf, []byte("\", \"line\":")...)
	buf = append(buf, []byte(fmt.Sprintf("%d", lg.GetLine()))...)
	buf = append(buf, []byte(", \"msg\":\"")...)
	buf = append(buf, []byte(lg.GetMsg())...)
	buf = append(buf, []byte("\"")...)
	tags := lg.GetTags()
	if len(tags) != 0 {
		buf = append(buf, []byte(", \"tags\":{")...)
		first := true
		for k, v := range tags {
			if !first {
				buf = append(buf, []byte(", ")...)
			}
			buf = append(buf, fmt.Sprintf("\"%s\":\"%s\"", k, v)...)
			first = false
		}
		buf = append(buf, []byte("}")...)
	}
	buf = append(buf, []byte("}")...)

	return string(buf)
}
