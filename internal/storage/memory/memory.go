package memory

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Memory struct {
	data map[string]*storage.Entity
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string]*storage.Entity)}
}

func (m *Memory) Exec(commands [][]byte) *protocol.ProtoValue {
	switch storage.Func(commands[0]) {
	case storage.SET:
		return m.set(commands[1:])
	case storage.GET:
		return m.get(commands[1:])
	default:
		return serrors.NewErrSyntaxIncorrect()
	}
}

func (m *Memory) set(args [][]byte) *protocol.ProtoValue {
	if len(args) != 2 {
		return serrors.NewErrSyntaxIncorrect()
	}
	m.data[string(args[0])] = storage.NewEntity(args[1])
	return serrors.NewOk()
}

func (m *Memory) get(args [][]byte) *protocol.ProtoValue {
	if len(args) != 1 {
		return serrors.NewErrSyntaxIncorrect()
	}
	if v, ok := m.data[string(args[0])]; ok {
		return serrors.NewBulkString(v.Value.([]byte))
	}
	return serrors.NewNilBulk()
}
