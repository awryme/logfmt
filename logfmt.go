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

type levelValues struct {
	Debug string
	Info  string
	Error string
}

// Params sets additional params for logger
type Params struct {
	// switch debug mode
	//default: enabled
	Debug bool

	// keys for standard fields
	// empty key remove the field from the line
	// defaults: lvl, msg, error
	LevelKey string
	MsgKey   string
	ErrorKey string

	// values for levels
	// type is not exported to simplify the api
	// empty value removes specific level from the line
	// defaults: debug, info, error
	Levels levelValues
}

// New create new logger with default params
// output should be valid, nil update allowed
func New(output io.Writer, update func(*Params)) Logger {
	params := &Params{
		Debug:    true,
		LevelKey: "lvl",
		MsgKey:   "msg",
		ErrorKey: "error",
		Levels: levelValues{
			Debug: "debug",
			Info:  "info",
			Error: "error",
		},
	}
	if update != nil {
		update(params)
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
	l.log(l.params.Levels.Info, msg, nil, kvs)
}

func (l loggerImpl) Debug(msg string, kvs ...Kv) {
	params := *l.params
	if params.Debug {
		l.log(params.Levels.Debug, msg, nil, kvs)
	}
}

func (l loggerImpl) Error(err error, msg string, kvs ...Kv) {
	if err != nil {
		l.log(l.params.Levels.Error, msg, err, kvs)
	}
}
