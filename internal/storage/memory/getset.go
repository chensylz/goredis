package memory

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

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
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	if v, ok := m.data[string(args[0])]; ok {
		return serrors.NewBulkString(v.Value.([]byte))
	}
	return serrors.NewNilBulk()
}
