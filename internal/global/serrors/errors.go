package serrors

import "github.com/chensylz/goredis/internal/protocol"

func NewErrProtocol() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR Protocol content error",
	}
}

func NewErrExec() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR Exec error",
	}
}

func NewErrUnknown() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR Unknown error",
	}
}

func NewOk() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.SimpleString,
		Value: "OK",
	}
}

func NewBulkString(value string) *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.BulkString,
		Value: value,
	}
}

func NewNilBulk() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.BulkString,
		Value: "",
	}
}

func NewErrSyntaxIncorrect() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR Syntax incorrect",
	}
}

func NewErrKeyNotFound() *protocol.ProtoValue {
	return &protocol.ProtoValue{
		Type:  protocol.Error,
		Value: "ERR Key not found",
	}
}
