package memory

import (
	"sync"
	"time"

	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Memory struct {
	data      map[string]*storage.Entity
	expireMap sync.Map
	sync.RWMutex
}

func NewMemory() *Memory {
	m := &Memory{data: make(map[string]*storage.Entity)}
	go func() {
		m.scanExpired()
		time.Sleep(100 * time.Millisecond)
	}()
	return m
}

func (m *Memory) Exec(commands *protocol.ProtoValue) *protocol.ProtoValue {
	value, ok := commands.Value.([]*protocol.ProtoValue)
	if !ok {
		return serrors.NewErrSyntaxIncorrect()
	}
	switch storage.Func(value[0].Value.(string)) {
	case storage.Set:
		return m.set(value[1:])
	case storage.Get:
		return m.get(value[1:])
	case storage.Expire:
		return m.expire(value[1:])
	case storage.Delete:
		return m.delete(value[1:])
	case storage.Ping:
		return m.ping()
	case storage.GetSet:
		return m.getSet(value[1:])
	default:
		return serrors.NewErrSyntaxIncorrect()
	}
}
