package encoder

import (
	"bytes"
	"io"
)

// Encoder encodes logfmt k=v pairs
//
// needs explicit WriteLine() to print log lines
type Encoder struct {
	writer    io.Writer
	lineHasKv bool
}

// New creates Encoder, that writes each k=v pair to writer directly
func New(writer io.Writer) *Encoder {
	return &Encoder{
		writer: writer,
	}
}

// WriteLine adds newline character
func (e *Encoder) WriteLine() error {
	e.lineHasKv = false
	_, err := e.writer.Write(bytesNewLine)
	return err
}

// WriteKvStr writes string k=v pair
//
// space is added between kvs implicitly
func (e *Encoder) WriteKvStr(key, value string) error {
	return e.WriteKv([]byte(key), []byte(value))
}

// WriteKv writes raw bytes k=v pair
//
// space is added between kvs implicitly
func (e *Encoder) WriteKv(key, value []byte) error {
	if key == nil {
		return nil
	}
	if e.lineHasKv {
		return e.writeAll(bytesSpace, formatKey(key), bytesEqual, formatValue(value))
	}
	e.lineHasKv = true
	return e.writeAll(formatKey(key), bytesEqual, formatValue(value))
}

func (e *Encoder) writeAll(bufs ...[]byte) error {
	for _, buf := range bufs {
		_, err := e.writer.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatKey(key []byte) []byte {
	key = escapeBuf(key)
	key = bytes.ReplaceAll(key, bytesSpace, []byte{'_'})
	return key
}

func formatValue(vOrig []byte) []byte {
	v := escapeBuf(vOrig)
	if len(vOrig) != len(v) || bytes.ContainsRune(v, symSpace) {
		var b bytes.Buffer
		b.Grow(len(v) + 2)
		b.WriteByte(symQuote)
		b.Write(v)
		b.WriteByte(symQuote)
		return b.Bytes()
	}
	return v
}

func escapeBuf(buf []byte) []byte {
	for _, r := range escapeSyms {
		buf = bytes.ReplaceAll(buf, []byte{r.sym}, []byte{symBackslash, r.replace})
	}
	return buf
}
