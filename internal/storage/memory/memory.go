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
	case storage.SET:
		return m.set(value[1:])
	case storage.GET:
		return m.get(value[1:])
	case storage.EXPIRE:
		return m.expire(value[1:])
	case storage.DELETE:
		return m.delete(value[1:])
	default:
		return serrors.NewErrSyntaxIncorrect()
	}
}
