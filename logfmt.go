package logfmt

import (
	"github.com/awryme/logfmt/encoder"
	"io"
)

// Logger logs msg on one of three levels
type Logger interface {
	// With adds values to the logger
	// Use logfmt.Delayed() to the value to be resolved at log time, not with time
	// Examples:
	// logger = logger.With(logfmt.Str("hello", "foo"), logfmt.Int("count", 3))
	// logger = logger.With(logfmt.Delayed(func() logfmt.Kv {
	//	 return logfmt.Str("time", time.Now().Format(time.RFC3339))
	// }))
	With(kvs ...Kv) Logger

	// Debug logs at debug level
	// Use `params.Debug = false` on logger creation to disable
	Debug(msg string, kvs ...Kv)

	// Info logs at info level
	Info(msg string, kvs ...Kv)

	// Error logs at error level
	// Ignored if err is nil
	Error(err error, msg string, kvs ...Kv)
}

// Params sets additional params for logger
type Params struct {
	// switch debug mode
	//default: enabled
	Debug bool

	// keys for standard fields
	// empty key remove the field from the line
	// defaults: "lvl", "msg", "error"
	LevelKey string
	MsgKey   string
	ErrorKey string

	// values for levels
	// empty value removes specific level from the line
	// defaults: "debug", "info", "error"
	LevelDebug string
	LevelInfo  string
	LevelError string
}

// New create new logger with default params
//
// output should be valid, nil options allowed
func New(output io.Writer, options func(*Params)) Logger {
	params := &Params{
		Debug: true,

		LevelKey: "lvl",
		MsgKey:   "msg",
		ErrorKey: "error",

		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelError: "error",
	}
	if options != nil {
		options(params)
	}
	return loggerImpl{
		output: output,
		params: params,
	}
}

// logger implementation
type loggerImpl struct {
	output io.Writer
	kvs    []Kv
	params *Params
}

func (l loggerImpl) log(level string, msg string, err error, kvs []Kv) {
	enc := encoder.NewLine()
	params := *l.params
	if params.LevelKey != "" && level != "" {
		enc.KvStr(params.LevelKey, level)
	}
	if params.MsgKey != "" {
		enc.KvStr(params.MsgKey, msg)
	}
	if err != nil && params.ErrorKey != "" {
		enc.KvStr(params.ErrorKey, err.Error())
	}
	for _, kv := range l.kvs {
		enc.Kv(kv.KvPair())
	}
	for _, kv := range kvs {
		enc.Kv(kv.KvPair())
	}
	_, _ = l.output.Write(enc.Bytes())
}

func (l loggerImpl) With(kvs ...Kv) Logger {
	l.kvs = append(l.kvs, kvs...)
	return l
}

func (l loggerImpl) Info(msg string, kvs ...Kv) {
	l.log(l.params.LevelInfo, msg, nil, kvs)
}

func (l loggerImpl) Debug(msg string, kvs ...Kv) {
	if l.params.Debug {
		l.log(l.params.LevelDebug, msg, nil, kvs)
	}
}

func (l loggerImpl) Error(err error, msg string, kvs ...Kv) {
	if err != nil {
		l.log(l.params.LevelError, msg, err, kvs)
	}
}
