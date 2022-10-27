package encoder

const (
	symNewLine   = '\n'
	symEqual     = '='
	symSpace     = ' '
	symQuote     = '"'
	symBackslash = '\\'
)

// need escape
var escapeSyms = []struct{ sym, replace byte }{
	// escape slash first
	// otherwise other slashes will be replaced
	{'\\', '\\'},
	{'\n', 'n'},
	{'\r', 'r'},
	{'\t', 't'},
	{'\b', 'b'},
	{'\f', 'f'},
	{'\v', 'v'},
	{'\000', '0'},
	{'"', '"'},
}

var (
	bytesNewLine = []byte{symNewLine}
	bytesEqual   = []byte{symEqual}
	bytesSpace   = []byte{symSpace}
)
