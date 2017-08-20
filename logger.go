/*
WoodLog is a logger downed from wood.

WoodLog is ridiculously simple structured and leveled logger  based on stdlib `log` package.
WoodLog is not lightning fast neither it doesn't produce extra complicated logs structure.
It does minimum what logger must do - L.O.G.S (!) with structure and levels.
*/
package woodlog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

// Low level logger
type Logger interface {
	// All required low level log methods
	Fatal(v ...interface{})
	Println(v ...interface{})
}

type formatter interface {
	formatSlots(map[string]interface{}) (bytes.Buffer, error)
}

// Base logger. Implements log structure format
type baseLog struct{}

// writeKV Writes structure key and casted value to a buffer.
// Returns formated dtructure field.
func (l *baseLog) writeKV(b *bytes.Buffer, k, v string) {
	b.WriteString(k)
	b.WriteString(": ")
	b.WriteString(v)
}

// formatSlots recursively formats log structure from slots.
// Returns buffer containing formatted log payload.
func (l *baseLog) formatSlots(slots map[string]interface{}) (buf bytes.Buffer, err error) {
	for k, v := range slots {
		switch t := v.(type) {
		case int:
			l.writeKV(&buf, k, strconv.Itoa(t))
		case bool:
			l.writeKV(&buf, k, strconv.FormatBool(t))
		case string:
			l.writeKV(&buf, k, t)
		default:
			_, f, l, _ := runtime.Caller(0)
			err = fmt.Errorf("%s %d: Wrong type of logging argument.", f, l)
		}
	}
	return
}

// Logger constructor.
type Log struct {
	formatter

	// Separate log.Logger instance for each log level.
	debug  Logger
	info   Logger
	error_ Logger
	trace  Logger
	fatal  Logger
}

// DEBUG level.
// Uses stdout as output file. Contains prefix "DEBUG", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) DEBUG(slots map[string]interface{}) (err error) {
	buf, err := l.formatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %s: Log DEBUG error.", err.Error(), f, l)
		return
	}
	l.debug.Println(buf.String())
	return
}

// INFO level.
// Uses stdout as output file. Contains prefix "INFO", microsecond resolution: 01:23:23.123123,
// final file name element and line number: d.go:23
func (l *Log) INFO(slots map[string]interface{}) (err error) {
	buf, err := l.formatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %s: Log INFO error.", err.Error(), f, l)
		return
	}
	l.info.Println(buf.String())
	return
}

// ERROR level.
// Uses stderr as output file. Contains prefix "ERROR", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) ERROR(slots map[string]interface{}) (err error) {
	buf, err := l.formatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %s: Log ERROR error.", err.Error(), f, l)
		return
	}
	l.error_.Println(buf.String())
	return
}

// TRACE level.
// Uses stdout as output file. Contains prefix "TRACE", microsecond resolution: 01:23:23.123123,
// final file name element and line number: d.go:23
func (l *Log) TRACE(slots map[string]interface{}) (err error) {
	buf, err := l.formatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %s: Log TRACE error.", err.Error(), f, l)
		return
	}
	l.trace.Println(buf.String())
	return
}

// FATAL level.
// Uses stderr as output file and exits with code 1 after logging. Contains prefix "FATAL", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) FATAL(slots map[string]interface{}) (err error) {
	buf, err := l.formatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %s: Log FATAL error.", err.Error(), f, l)
		return
	}
	l.fatal.Fatal(buf.String())
	return
}

// New logger
func New() *Log {
	return &Log{
		debug:  log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
		info:   log.New(os.Stdout, "INFO: ", log.Lmicroseconds|log.Lshortfile),
		error_: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
		trace:  log.New(os.Stdout, "TRACE: ", log.Lmicroseconds|log.Lshortfile),
		fatal:  log.New(os.Stderr, "FATAL", log.Ldate|log.Lmicroseconds|log.Llongfile),
	}
}
