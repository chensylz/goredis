package memory

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

func (m *Memory) set(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) != 2 {
		return serrors.NewErrSyntaxIncorrect()
	}
	value := args[1].Value.(string)
	m.data[args[0].Value.(string)] = storage.NewEntity(value, int64(len(value)))
	return serrors.NewOk()
}

func (m *Memory) get(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) != 1 {
		return serrors.NewErrSyntaxIncorrect()
	}
	m.Lock()
	defer m.Unlock()
	key := args[0].Value.(string)
	v, ok := m.data[key]
	if !ok {
		return serrors.NewNilBulk()
	}
	v.Hit++
	m.data[key] = v
	return serrors.NewBulkString(v.Value.(string))
}
