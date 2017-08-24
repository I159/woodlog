package woodlog

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
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
				t.Error("Incorrect formatting: %s", b.String())
			}
		})
	}
}

func Test_baseLog_formatSlots(t *testing.T) {
	var bufInt, bufBool, bufStr bytes.Buffer
	bufBool.WriteString("k: false")
	bufInt.WriteString("k: 1")
	bufStr.WriteString("k: a")
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		l       *baseLog
		args    args
		wantBuf bytes.Buffer
		wantErr bool
	}{
		{
			name:    "Wrong argument type",
			wantErr: true,
			args:    args{map[string]interface{}{"k": int64(1)}},
		},
		{
			name:    "int value",
			wantBuf: bufInt,
			args:    args{map[string]interface{}{"k": 1}},
		},
		{
			name:    "bool value",
			wantBuf: bufBool,
			args:    args{map[string]interface{}{"k": false}},
		},
		{
			name:    "string value",
			wantBuf: bufStr,
			args:    args{map[string]interface{}{"k": "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &baseLog{}
			gotBuf, err := l.formatSlots(tt.args.slots)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseLog.formatSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("baseLog.formatSlots() = %v, want %v", gotBuf, tt.wantBuf)
			}
		})
	}
}

type mockFailFormatter struct{}

func (m *mockFailFormatter) formatSlots(map[string]interface{}) (b bytes.Buffer, err error) {
	err = errors.New("")
	return
}

func TestLog_DEBUG(t *testing.T) {
	var kvBuf bytes.Buffer
	var b bytes.Buffer

	kvBuf.WriteString("")

	type fields struct {
		formatter formatter
		debug     Logger
		info      Logger
		error_    Logger
		trace     Logger
		fatal     Logger
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf *bytes.Buffer
	}{
		{
			name: "Key value",
			fields: fields{
				formatter: new(baseLog),
				debug:     log.New(&b, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Llongfile),
			},
			args: args{map[string]interface{}{"k": "v"}},
			wantBuf: &kvBuf,
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
				error_:    tt.fields.error_,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.DEBUG(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.DEBUG() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantBuf != nil &&  strings.Compare(b.String(), tt.wantBuf.String()) != 0 {
				t.Errorf("Logger output: %v; want: %v", b.String(), tt.wantBuf.String())
			}
		})
	}
}

func TestLog_INFO(t *testing.T) {
	type fields struct {
		formatter formatter
		debug     Logger
		info      Logger
		error_    Logger
		trace     Logger
		fatal     Logger
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				error_:    tt.fields.error_,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.INFO(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.INFO() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLog_ERROR(t *testing.T) {
	type fields struct {
		formatter formatter
		debug     Logger
		info      Logger
		error_    Logger
		trace     Logger
		fatal     Logger
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				error_:    tt.fields.error_,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.ERROR(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.ERROR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLog_TRACE(t *testing.T) {
	type fields struct {
		formatter formatter
		debug     Logger
		info      Logger
		error_    Logger
		trace     Logger
		fatal     Logger
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				error_:    tt.fields.error_,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.TRACE(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.TRACE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLog_FATAL(t *testing.T) {
	type fields struct {
		formatter formatter
		debug     Logger
		info      Logger
		error_    Logger
		trace     Logger
		fatal     Logger
	}
	type args struct {
		slots map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Log{
				formatter: tt.fields.formatter,
				debug:     tt.fields.debug,
				info:      tt.fields.info,
				error_:    tt.fields.error_,
				trace:     tt.fields.trace,
				fatal:     tt.fields.fatal,
			}
			if err := l.FATAL(tt.args.slots); (err != nil) != tt.wantErr {
				t.Errorf("Log.FATAL() error = %v, wantErr %v", err, tt.wantErr)
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
