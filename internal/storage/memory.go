package storage

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
)

type Memory struct {
	data map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string]interface{})}
}

func (m *Memory) Exec(commands [][]byte) *protocol.ProtoValue {
	switch string(commands[0]) {
	case "SET":
		return m.set(commands[1:])
	case "GET":
		return m.get(commands[1:])
	default:
		return serrors.NewErrSyntaxIncorrect()
	}
}

func (m *Memory) set(args [][]byte) *protocol.ProtoValue {
	if len(args) != 2 {
		return serrors.NewErrSyntaxIncorrect()
	}
	m.data[string(args[0])] = string(args[1])
	return serrors.NewOk()
}

func (m *Memory) get(args [][]byte) *protocol.ProtoValue {
	if len(args) != 1 {
		return serrors.NewErrSyntaxIncorrect()
	}
	if v, ok := m.data[string(args[0])]; ok {
		return serrors.NewBulkString(v.(string))
	}
	return serrors.NewNilBulk()
}
