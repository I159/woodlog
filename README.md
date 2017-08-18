## **WoodLog**

WoodLog is a logger downed from wood.

WoodLog is ridiculously simple structured and leveled logger based on
stdlib `log` package. WoodLog is not lightning fast neither it doesn't
produce extra complicated logs structure. It does minimum what logger
must do - L.O.G.S (!) with structure and levels.

### **TYPES**

`type BaseLog struct{}`
    Base logger. Implements log structure format

    type Log struct {
        BaseLog
        // contains filtered or unexported fields
     }
Logger implementation

`func New(level string) (logger *Log, err error)`
New logger

`func (l *Log) DEBUG(slots map[string]interface{}) (err error)`
    DEBUG level. Uses stdout as output file. Contains prefix "DEBUG", the
    date in the local time zone: 2009/01/23, microsecond resolution:
    01:23:23.123123, full file name and line number: /a/b/c/d.go:23.

`func (l *Log) ERROR(slots map[string]interface{}) (err error)`
    ERROR level. Uses stderr as output file. Contains prefix "ERROR", the
    date in the local time zone: 2009/01/23, microsecond resolution:
    01:23:23.123123, full file name and line number: /a/b/c/d.go:23.

`func (l *Log) FATAL(slots map[string]interface{}) (err error)`
    FATAL level. Uses stderr as output file and exits with code 1 after
    logging. Contains prefix "FATAL", the date in the local time zone:
    2009/01/23, microsecond resolution: 01:23:23.123123, full file name and
    line number: /a/b/c/d.go:23.

`func (l *Log) INFO(slots map[string]interface{}) (err error)`
    INFO level. Uses stdout as output file. Contains prefix "INFO",
    microsecond resolution: 01:23:23.123123, final file name element and
    line number: d.go:23.

`func (l *Log) TRACE(slots map[string]interface{}) (err error)`
    TRACE level. Uses stdout as output file. Contains prefix "TRACE",
    microsecond resolution: 01:23:23.123123, final file name element and
    line number: d.go:23.

    type Logger interface {
        DEBUG(map[string]interface{}) error
        INFO(map[string]interface{}) error
        ERROR(map[string]interface{}) error
        FATAL(map[string]interface{}) error
        TRACE(map[string]interface{}) error
    }
Public interface to extend existing functional.

