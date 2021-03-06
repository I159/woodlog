package woodlog

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func Test_baseLog_writeKV(t *testing.T) {
	var b bytes.Buffer

	type args struct {
		b *bytes.Buffer
		k string
		v string
	}
	tests := []struct {
		name string
		l    *baseLog
		args args
	}{
		{
			"Normal key & value",
			new(baseLog),
			args{
				&b,
				"k",
				"v",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &baseLog{}
			l.writeKV(tt.args.b, tt.args.k, tt.args.v)
			defer b.Reset()

			if strings.Compare(b.String(), "k: v") != 0 {
				t.Errorf("Incorrect formatting: %s", b.String())
			}
		})
	}
}

func Test_baseLogFormatSlots(t *testing.T) {
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		l       *baseLog
		args    args
		wantBuf *bytes.Buffer
		wantErr bool
	}{
		{
			name:    "Wrong argument type",
			wantErr: true,
			args:    args{map[string]interface{}{"k": int64(1)}},
			wantBuf: new(bytes.Buffer),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &baseLog{}
			gotBuf, err := l.FormatSlots(tt.args.slots)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseLog.FormatSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("baseLog.FormatSlots() = %v, want %v", gotBuf, tt.wantBuf)
			}
		})
	}
}

type mockFailFormatter struct{}

func (m *mockFailFormatter) FormatSlots(map[string]interface{}) (b *bytes.Buffer, err error) {
	b = bytes.NewBuffer([]byte{})
	err = errors.New("")
	return
}

func TestLog_DEBUG(t *testing.T) {
	var kvBuf bytes.Buffer
	var b bytes.Buffer

	kvBuf.WriteString("")

	type fields struct {
		formatter formatter
		debug     LoggerLevel
		info      LoggerLevel
		eRRor     LoggerLevel
		trace     LoggerLevel
		fatal     LoggerLevel
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantPattern *regexp.Regexp
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				debug:     log.New(&b, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
			},
			args:        args{map[string]interface{}{"k": "v"}},
			wantPattern: regexp.MustCompile("DEBUG: [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6} [a-zA-Z0-9/_\\.]*:[0-9]+: [a-zA-Z0-9: ]*"),
		},
		{
			name: "Format fail",
			fields: fields{
				formatter: new(mockFailFormatter),
				debug:     log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
			},
			args:    args{map[string]interface{}{"k": "v"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				eRRor:     tt.fields.eRRor,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.DEBUG(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.DEBUG() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantPattern != nil && tt.wantPattern.MatchString(b.String()) != true {
				t.Errorf("Logger output: %v; want pattern: %v", b.String(), tt.wantPattern.String())
			}
		})
	}
}

func TestLog_INFO(t *testing.T) {
	var b bytes.Buffer

	type fields struct {
		formatter formatter
		debug     LoggerLevel
		info      LoggerLevel
		eRRor     LoggerLevel
		trace     LoggerLevel
		fatal     LoggerLevel
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantPattern *regexp.Regexp
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				info:      newINFO(&b),
			},
			args:        args{map[string]interface{}{"k": "v"}},
			wantPattern: regexp.MustCompile("INFO: [0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6} [a-zA-Z0-9/_\\.]*:[0-9]+: [a-zA-Z0-9: ]+"),
		},
		{
			name: "Empty payload",
			fields: fields{
				formatter: new(baseLog),
				info:      newINFO(&b),
			},
			wantErr: true,
		},
		{
			name: "Format fail",
			fields: fields{
				formatter: new(mockFailFormatter),
				info:      newINFO(&b),
			},
			args:    args{map[string]interface{}{"k": "v"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				eRRor:     tt.fields.eRRor,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.INFO(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.INFO() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantPattern != nil && tt.wantPattern.MatchString(b.String()) != true {
				t.Errorf("Logger output: %v; want pattern: %v", b.String(), tt.wantPattern.String())
			}
		})
	}
}

func TestLog_ERROR(t *testing.T) {
	var b bytes.Buffer

	type fields struct {
		formatter formatter
		debug     LoggerLevel
		info      LoggerLevel
		eRRor     LoggerLevel
		trace     LoggerLevel
		fatal     LoggerLevel
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantPattern *regexp.Regexp
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				eRRor:     newERROR(&b),
			},
			args:        args{map[string]interface{}{"k": "v"}},
			wantPattern: regexp.MustCompile("ERROR: [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6} [a-zA-Z0-9/_\\.]*:[0-9]+: [a-zA-Z0-9: ]+"),
		},
		{
			name: "Empty payload",
			fields: fields{
				formatter: new(baseLog),
				eRRor:     newERROR(&b),
			},
			wantErr: true,
		},
		{
			name: "Format fail",
			fields: fields{
				formatter: new(mockFailFormatter),
				eRRor:     newERROR(&b),
			},
			args:    args{map[string]interface{}{"k": "v"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				eRRor:     tt.fields.eRRor,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.ERROR(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.ERROR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantPattern != nil && tt.wantPattern.MatchString(b.String()) != true {
				t.Errorf("Logger output: %v; want pattern: %v", b.String(), tt.wantPattern.String())
			}
		})
	}
}

func TestLog_TRACE(t *testing.T) {
	var b bytes.Buffer

	type fields struct {
		formatter formatter
		debug     LoggerLevel
		info      LoggerLevel
		eRRor     LoggerLevel
		trace     LoggerLevel
		fatal     LoggerLevel
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantPattern *regexp.Regexp
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				trace:     newTRACE(&b),
			},
			args:        args{map[string]interface{}{"k": "v"}},
			wantPattern: regexp.MustCompile("TRACE: [0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6} [a-zA-Z0-9/_\\.]*:[0-9]+: [a-zA-Z0-9: ]+"),
		},
		{
			name: "Empty payload",
			fields: fields{
				formatter: new(baseLog),
				trace:     newTRACE(&b),
			},
			wantErr: true,
		},
		{
			name: "Format fail",
			fields: fields{
				formatter: new(mockFailFormatter),
				trace:     newTRACE(&b),
			},
			args:    args{map[string]interface{}{"k": "v"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				eRRor:     tt.fields.eRRor,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.TRACE(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.TRACE() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantPattern != nil && tt.wantPattern.MatchString(b.String()) != true {
				t.Errorf("Logger output: %v; want pattern: %v", b.String(), tt.wantPattern.String())
			}
		})
	}
}

type mockLogger struct {
	logger *log.Logger
}

func (m *mockLogger) Fatal(v ...interface{}) {
	m.logger.Println(v...)
}
func (m *mockLogger) Println(v ...interface{}) {
	m.logger.Println(v...)
}

func TestLog_FATAL(t *testing.T) {
	var b bytes.Buffer

	type fields struct {
		formatter formatter
		debug     LoggerLevel
		info      LoggerLevel
		eRRor     LoggerLevel
		trace     LoggerLevel
		fatal     LoggerLevel
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantPattern *regexp.Regexp
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				fatal:     &mockLogger{log.New(&b, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Llongfile)},
			},
			args:        args{map[string]interface{}{"k": "v"}},
			wantPattern: regexp.MustCompile("FATAL: [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6} [a-zA-Z0-9/_\\.]*:[0-9]+: [a-zA-Z0-9: ]+"),
		},
		{
			name: "Empty payload",
			fields: fields{
				formatter: new(baseLog),
				fatal:     &mockLogger{log.New(&b, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Llongfile)},
			},
			wantErr: true,
		},
		{
			name: "Format fail",
			fields: fields{
				formatter: new(mockFailFormatter),
				fatal:     &mockLogger{log.New(&b, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Llongfile)},
			},
			args:    args{map[string]interface{}{"k": "v"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				eRRor:     tt.fields.eRRor,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.FATAL(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.FATAL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantPattern != nil && tt.wantPattern.MatchString(b.String()) != true {
				t.Errorf("Logger output: %v; want pattern: %v", b.String(), tt.wantPattern.String())
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Log
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
