package encoder

import "bytes"

// Line represents single log line
// Adds newline character implicitly
type Line struct {
	line    *bytes.Buffer
	encoder *Encoder
}

// NewLine create Line with new buffer
func NewLine() *Line {
	buf := bytes.NewBuffer(nil)
	return &Line{
		line:    buf,
		encoder: New(buf),
	}
}

// KvStr writes string k=v pair to buffer
func (l *Line) KvStr(key, value string) {
	_ = l.encoder.WriteKvStr(key, value)
}

// Kv writes raw bytes k=v pair to buffer
func (l *Line) Kv(key, value []byte) {
	_ = l.encoder.WriteKv(key, value)
}

// Bytes return a log line of k=v pairs with spaces between them and a newline character
func (l *Line) Bytes() []byte {
	_ = l.encoder.WriteLine()
	return l.line.Bytes()
}
