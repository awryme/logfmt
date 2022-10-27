package logfmt

import (
	"fmt"
	"strconv"
)

type Kv struct {
	Key, Value []byte
	Delayed    func() Kv
}

func (kv Kv) KvPair() ([]byte, []byte) {
	if kv.Delayed != nil {
		return kv.Delayed().KvPair()
	}
	return kv.Key, kv.Value
}

// Delayed logs Kv value, resolved at log time
//
// can be used to add values like timestamp/caller for all logs using logger.With()
func Delayed(f func() Kv) Kv {
	return Kv{
		Delayed: f,
	}
}

// Bytes logs raw bytes as is
func Bytes(key string, value []byte) Kv {
	return Kv{
		Key:   []byte(key),
		Value: value,
	}
}

// Str logs string value
func Str(key string, value string) Kv {
	return Kv{
		Key:   []byte(key),
		Value: []byte(value),
	}
}

// Stringer logs value.String()
//
// nil value is ignored
func Stringer(key string, value fmt.Stringer) Kv {
	if value == nil {
		return Kv{}
	}
	return Str(key, value.String())
}

// Int logs int values
//
// cast to int64 if necessary
func Int(key string, value int64) Kv {
	return Kv{
		Key:   []byte(key),
		Value: strconv.AppendInt(nil, value, 10),
	}
}

// Uint logs uint values
//
// cast to uint64 if necessary
func Uint(key string, value uint64) Kv {
	return Kv{
		Key:   []byte(key),
		Value: strconv.AppendUint(nil, value, 10),
	}
}

// Float logs float values without a precision limit
//
// cast to float64 if necessary
func Float(key string, value float64) Kv {
	return Kv{
		Key:   []byte(key),
		Value: strconv.AppendFloat(make([]byte, 0, 24), value, 'f', -1, 64),
	}
}

// Bool logs bool values "true"/"false"
func Bool(key string, value bool) Kv {
	return Kv{
		Key:   []byte(key),
		Value: strconv.AppendBool(nil, value),
	}
}

// Any logs any value with fmt.Sprint
//
// Added for convenience
func Any(key string, value interface{}) Kv {
	return Kv{
		Key:   []byte(key),
		Value: []byte(fmt.Sprint(value)),
	}
}
