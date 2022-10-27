# logfmt

[![Go Reference](https://pkg.go.dev/badge/github.com/awryme/logfmt.svg)](https://pkg.go.dev/github.com/awryme/logfmt)

Package logfmt provide a simple logger in [logfmt format](https://brandur.org/logfmt)

+ Spaces in keys are replaced with underscores, empty keys are ignored
+ Values with spaces are wrapped in double quotes
+ Tags are not supported

---

Logger uses simple api

Kv is a simple struct of []byte values, you can create your own

Helper functions are defined
```go
type Logger interface {
	With(kvs ...Kv) Logger
	Debug(msg string, kvs ...Kv)
	Info(msg string, kvs ...Kv)
	Error(err error, msg string, kvs ...Kv)
}

type Kv struct {
	Key, Value []byte
	Delayed    func() Kv
}

func Bytes(key string, value []byte) Kv
func Str(key string, value string) Kv
func Stringer(key string, value fmt.Stringer) Kv
func Int(key string, value int64) Kv
func Uint(key string, value uint64) Kv
func Float(key string, value float64) Kv
func Bool(key string, value bool) Kv
func Any(key string, value interface{}) Kv
func Delayed(f func() Kv) Kv
```

Example:
```go
package main

import (
	"fmt"
	"github.com/awryme/logfmt"
	"os"
	"sync/atomic"
)

var count int32 = 10

func main() {
	logger := logfmt.New(os.Stderr, func(params *logfmt.Params) {
		// Disable debug (default: enabled)
		params.Debug = false

		// Do not print levels
		params.LevelKey = ""
	})

	// delayed until log time
	// spaces converted to underscore
	// value needs cast to int64
	logger = logger.With(logfmt.Delayed(func() logfmt.Kv {
		value := atomic.AddInt32(&count, 2)
		return logfmt.Int("global count", int64(value))
	}))

	// static values
	logger = logger.With(logfmt.Str("appname", "test"), logfmt.Float("pi value", 3.14))

	// not printed, debug is disabled
	// count is not increased
	logger.Debug("some debug")

	// nil error logs are ignored
	logger.Error(nil, "failure or not?")

	// "local" kv are printed after "global" ones
	logger.Info("starting test app", logfmt.Int("local-count", 3))

	err := fmt.Errorf("failed everything")
	logger.Error(err, "failure on start")
}

// result stderr:
// msg="starting test app" global_count=12 appname=test pi_value=3.14 local-count=3
// msg="failure on start" error="failed everything" global_count=14 appname=test pi_value=3.14

```
