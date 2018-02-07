/*
Package woodlog is a logger downed from wood.

WoodLog is ridiculously simple structured and leveled logger  based on stdlib `log` package.
WoodLog is not lightning fast neither it doesn't produce extra complicated logs structure.
It does minimum what logger must do - L.O.G.S (!) with structure and levels.
## Motivation

There is plenty of loggers written on Go. All of it extra complicated
lightning fast mega-flexible huge like a space ships libraries. It's all too mach
if you are building a small tool and want to know what happens
inside of it. Skipping low level optimization, contexts and flexibility WoodLog
implements simple as a piece of wood structured and
leveled logger.

## Installation

Probably you have Go installed and configured. So just do:

    `go get github.com/I159/woodlog`

Then you can import it and use it as described below.

## Get started

Log everything you need to expose out of an application using simple system of levels
and map to pass a key-value structure of logged arguments to the logger instance.
Let's say you have an http server and you want to bring some transparency to it:

	package main

	import (
		"fmt"
		"html"
		"net/http"

		"github.com/I159/woodlog"
	)

	var logger *woodlog.Log = woodlog.New()

	func importantFunc() (int , error) {
		return 0, nil
	}

	func handler(w http.ResponseWriter, r *http.Request) {
		logger.TRACE(map[string]interface{}{"Request state": "Received"})
		logger.INFO(map[string]interface{}{"Requested url": "/bar"})

		if num, err := importantFunc(); err == nil {
			logger.DEBUG(map[string]interface{}{"Incoming number": num})
			fmt.Fprintf(w, "Hello, %q #%d", html.EscapeString(r.URL.Path), num)
		} else {
			logger.ERROR(map[string]interface{}{"Error": err, "Occured in": "importantFunc"})
		}

		logger.TRACE(map[string]interface{}{"Request state": "Responded"})
	}

	func main() {
		http.HandleFunc("/bar", handler)
		http.ListenAndServe(":8080", nil)
	}

## Advanced usage and custom loggers

If existing granularity of levels is not enough you can implement your own logger based on `formatter` interface.
See `Logger` interface documentation for details.

	package main                                                                                    

	import (
		"fmt"
		"log"
		"os"
		"runtime"

		"github.com/I159/woodlog"
	)

	type CustomLog struct {
		*woodlog.Log
		panic woodlog.LoggerLevel
	}

	func (l *CustomLog) PANIC(slots map[string]interface{}) (err error) {
		buf, err := l.FormatSlots(slots)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			err = fmt.Errorf("%s\n%s %d: Log DEBUG error.", err.Error(), f, l)
			return
		}
		defer panic(buf.String())
		l.panic.Println()
		return
	}

	func newCustomLogger() *CustomLog {
		return &CustomLog{
			woodlog.New(),
			log.New(os.Stderr, "PANIC: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
		}
	}

	func main() {
		logger := newCustomLogger()
		logger.PANIC(map[string]interface{}{"Cause of a panic": "Zombie! Zombies are everywhere!"})
	}
*/
package woodlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
)

// LoggerLevel interface contains only the must methods for implementetion
// of logger levels.
type LoggerLevel interface {
	// All required low level log methods
	Fatal(v ...interface{})
	Println(v ...interface{})
}

type formatter interface {
	FormatSlots(map[string]interface{}) (*bytes.Buffer, error)
}

// Base logger. Implements log structure format
type baseLog struct{}

// writeKV
// Writes a passed key-value pairs to a *bytes.Buffer as a formatted string.
func (l *baseLog) writeKV(b *bytes.Buffer, k, v string) {
	b.WriteString(k)
	b.WriteString(": ")
	b.WriteString(v)
}

// FormatSlots recursively formats log structure from slots.
// Returns buffer containing formatted log payload.
func (l *baseLog) FormatSlots(slots map[string]interface{}) (buf *bytes.Buffer, err error) {
	if len(slots) == 0 {
		err = errors.New("empty logging is not allowed")
	}

	buf = new(bytes.Buffer)
	for k, v := range slots {
		switch t := v.(type) {
		case int:
			l.writeKV(buf, k, strconv.Itoa(t))
		case bool:
			l.writeKV(buf, k, strconv.FormatBool(t))
		case string:
			l.writeKV(buf, k, t)
		default:
			_, f, l, _ := runtime.Caller(0)
			err = fmt.Errorf("%s %d: Wrong type of logging argument", f, l)
		}
	}
	return
}

// Log is a composition of loggers of different levels. Each logger satisfies `Logger` interface
// and could be replaced with a custom object satisfying the interface.
type Log struct {
	formatter

	debug LoggerLevel
	info  LoggerLevel
	eRRor LoggerLevel
	trace LoggerLevel
	fatal LoggerLevel
}


// TODO: refactor all the log methods wuth this generalization function
func (l *Log) fillBuffer(tmpl string, slots map[string]interface{}) (buf fmt.Stringer, err error) {
	buf, err = l.FormatSlots(slots)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		err = fmt.Errorf("%s\n%s %d: Log DEBUG error", err.Error(), f, l)
	}
	return
}


// DEBUG level.
// Uses stdout as output file. Contains prefix "DEBUG", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) DEBUG(slots map[string]interface{}) (err error) {
	buf, err := l.fillBuffer("%s\n%s %d: Log DEBUG error", slots)
	if err != nil {
		return
	}
	l.debug.Println(buf.String())
	return
}

// INFO level.
// Uses stdout as output file. Contains prefix "INFO", microsecond resolution: 01:23:23.123123,
// final file name element and line number: d.go:23
func (l *Log) INFO(slots map[string]interface{}) (err error) {
	buf, err := l.fillBuffer("%s\n%s %d: Log INFO error", slots)
	if err != nil {
		return
	}
	l.info.Println(buf.String())
	return
}

// ERROR level.
// Uses stderr as output file. Contains prefix "ERROR", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) ERROR(slots map[string]interface{}) (err error) {
	buf, err := l.fillBuffer("%s\n%s %d: Log ERROR error", slots)
	if err != nil {
		return
	}
	l.eRRor.Println(buf.String())
	return
}

// TRACE level.
// Uses stdout as output file. Contains prefix "TRACE", microsecond resolution: 01:23:23.123123,
// final file name element and line number: d.go:23
func (l *Log) TRACE(slots map[string]interface{}) (err error) {
	buf, err := l.fillBuffer("%s\n%s %d: Log TRACE error", slots)
	if err != nil {
		return
	}
	l.trace.Println(buf.String())
	return
}

// FATAL level.
// Uses stderr as output file and exits with code 1 after logging. Contains prefix "FATAL", the date in the local time zone: 2009/01/23,
// microsecond resolution: 01:23:23.123123, full file name and line number: /a/b/c/d.go:23.
func (l *Log) FATAL(slots map[string]interface{}) (err error) {
	buf, err := l.fillBuffer("%s\n%s %d: Log FATAL error", slots)
	if err != nil {
		return
	}
	l.fatal.Fatal(buf.String())
	return
}

func newDEBUG(wr io.Writer) LoggerLevel {
	return log.New(wr, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

func newINFO(wr io.Writer) LoggerLevel {
	return log.New(wr, "INFO: ", log.Lmicroseconds|log.Lshortfile)
}

func newERROR(wr io.Writer) LoggerLevel {
	return log.New(wr, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

func newTRACE(wr io.Writer) LoggerLevel {
	return log.New(wr, "TRACE: ", log.Lmicroseconds|log.Lshortfile)
}

func newFATAL(wr io.Writer) LoggerLevel {
	return log.New(wr, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

// New logger
func New() *Log {
	return &Log{
		formatter: new(baseLog),
		debug:     newDEBUG(os.Stdout),
		info:      newINFO(os.Stdout),
		eRRor:     newERROR(os.Stderr),
		trace:     newTRACE(os.Stdout),
		fatal:     newFATAL(os.Stderr),
	}
}
