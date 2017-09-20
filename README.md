[![Build Status](https://travis-ci.org/I159/woodlog.svg?branch=master)](https://travis-ci.org/I159/woodlog)
[![Code Climate](https://codeclimate.com/github/I159/woodlog/badges/gpa.svg)](https://codeclimate.com/github/I159/woodlog)
[![Test Coverage](https://codeclimate.com/github/I159/woodlog/badges/coverage.svg)](https://codeclimate.com/github/I159/woodlog/coverage)
# woodlog

Package woodlog is a logger downed from wood.

WoodLog is ridiculously simple structured and leveled logger based on
stdlib `log` package. WoodLog is not lightning fast neither it doesn't
produce extra complicated logs structure. It does minimum what logger
must do - L.O.G.S (!) with structure and levels.

## Motivation

There is plenty of loggers written on Go. All of it are extra complicated
lightning fast mega-flexible huge like a space ships libraries. It's all
too much if you are building a small tool and want to know what happens
inside of it. Skipping low level optimization, contexts and flexibility
WoodLog implements simple as a piece of wood structured and leveled
logger.

## Installation

Probably you have Go installed and configured. So just do:

	go get github.com/I159/woodlog

Then you can import it and use it as described below.

## Get started

Log everything you need to expose out of an application using simple
system of levels and map to pass a key-value structure of logged
arguments to the logger instance. Let's say you have an http server and
you want to bring some transparency to it:

	package main

	import (
		"fmt"
		"html"
		"net/http"

		"github.com/I159/woodlog"
	)

	var logger *woodlog.Logger = woodlog.New()

	func handler(w http.ResponseWriter, r *http.Request) {
		logger.TRACE(map[string]interface{}{"Request state": "received"})
		logger.INFO(map[string]interface{}{"Requested url": "/bar"})

		if num, err := importantFunc(); err == nil {
			logger.DEBUG(map[string]interface{}{"Incoming word number": num})
			fmt.Fprintf(w, "Hello, %q #%d", html.EscapeString(r.URL.Path), num)
		} else {
			logger.ERROR(map[string]interface{}{"Error": err, "Occured in": "importantFunc"})
		}

		logger.TRACE(map[string]interface{}{"Request state": "jesponded"})
	}

	func main() {
		http.HandleFunc("/bar", handler)
		logger.Fatal(map[string]interface{}{"Server stopped due to": http.ListenAndServe(":8080", nil)})
	}

## Advanced usage and custom loggers

If existing granularity of levels is not enough you can implement your
own logger based on `formatter` interface. See `Logger` interface
documentation for details.

	type CustomLog struct {
		*Log
		panic *Logger
	}

	func (l *CustomLog) PANIC(slots map[string]interface{}) (err error) {
		buf, err := l.FormatSlots(slots)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			err = fmt.Errorf("%s\n%s %d: Log DEBUG error.", err.Error(), f, l)
			return
		}
		defer panic(buf.String())
		l.panic.Println(buf.String())
		return
	}

	func newCustomLogger() *CustomLog {
		return &CustomLog{
			woodlog.New(),
			log.New(wr, "PANIC: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
		}
	}

## API reference

    type Log struct {
        // contains filtered or unexported fields
    }
Log is a composition of loggers of different levels. Each logger
satisfies `LoggerLevel` interface and could be replaced with a custom object
satisfying the interface.

    func New() *Log
New logger constructor function.

    func (l *Log) DEBUG(slots map[string]interface{}) (err error)
DEBUG level. Uses stdout as output file. Contains prefix "DEBUG", the
date in the local time zone: 2009/01/23, microsecond resolution:
01:23:23.123123, full file name and line number: /a/b/c/d.go:23.

func (l *Log) ERROR(slots map[string]interface{}) (err error)
    ERROR level. Uses stderr as output file. Contains prefix "ERROR", the
    date in the local time zone: 2009/01/23, microsecond resolution:
    01:23:23.123123, full file name and line number: /a/b/c/d.go:23.

    func (l *Log) FATAL(slots map[string]interface{}) (err error)
FATAL level. Uses stderr as output file and exits with code 1 after
logging. Contains prefix "FATAL", the date in the local time zone:
2009/01/23, microsecond resolution: 01:23:23.123123, full file name and
line number: /a/b/c/d.go:23.

    func (l *Log) INFO(slots map[string]interface{}) (err error)
INFO level. Uses stdout as output file. Contains prefix "INFO",
microsecond resolution: 01:23:23.123123, final file name element and
line number: d.go:23

    func (l *Log) TRACE(slots map[string]interface{}) (err error)
TRACE level. Uses stdout as output file. Contains prefix "TRACE",
microsecond resolution: 01:23:23.123123, final file name element and
line number: d.go:23

    type LoggerLevel interface {
        // All required low level log methods
        Fatal(v ...interface{})
        Println(v ...interface{})
    }
LoggerLevel interface contains only the must methods for implementation
of logger levels.

## Contribution

Desired level of test coverage is 80%. Please don't commit changes with lower level of coverage.
Commit messages should be in this format:

    Topic which describes changes in plain text

    Detailed description after a blank line not more than
    72 chracters per line.
	    func Code(allowed bool) string {
            // Comment
	    }
    Topic shouldn't be lesser than 15 characters.
    Description ends with Closes|Resolves: #<github issue num>
    Closes: #3

Commit message style controlled with client side hook which you can find in
.template directory of the repo.

Don't leave `println` or `fmt.Print*` statements. Seriously! This is the logger =) Controlled with hook too.
Keep the hooks in `.git/hook` directory, it made for smoother contribution.

Keep unfinished changes in a feature branches and make pull requests only when you believe it closes an issue.

## Plans

TRACE level ability to wrap functions and return time report.
Make hooks for test coverage.
