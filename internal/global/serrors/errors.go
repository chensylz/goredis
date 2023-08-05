package serrors

const symbol = "\r\n"

var (
	ErrProtocol        = []byte("-ERR Protocol content error" + symbol)
	ErrExec            = []byte("-ERR Exec error" + symbol)
	ErrSyntaxIncorrect = []byte("-ERR Syntax incorrect" + symbol)
	Ok                 = []byte("+OK" + symbol)
	NilBulk            = []byte("$-1" + symbol)
)
