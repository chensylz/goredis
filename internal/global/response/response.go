package response

import "github.com/chensylz/goredis/internal/protocol"

const (
	SyntaxIncorrect = "Syntax incorrect"
)

var (
	ProtocolErr        = NewErr("Protocol content error")
	UnknownErr         = NewErr("Unknown error")
	SyntaxIncorrectErr = NewErr(SyntaxIncorrect)

	Ok      = NewSimpleString("Ok")
	NilBulk = NewBulkString("")
	One     = NewInter(1)
)

func NewErr(message string) *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR" + message,
	}
}

func NewSimpleString(value string) *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.SimpleString,
		Value: value,
	}
}

func NewInter(value int64) *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Integer,
		Value: value,
	}
}

func NewBulkString(value string) *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.BulkString,
		Value: value,
	}
}
