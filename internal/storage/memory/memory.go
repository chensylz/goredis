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

func (m *Memory) Exec(commands [][]byte) *protocol.ProtoValue {
	switch storage.Func(commands[0]) {
	case storage.SET:
		return m.set(commands[1:])
	case storage.GET:
		return m.get(commands[1:])
	case storage.EXPIRE:
		return m.expire(commands[1:])
	default:
		return serrors.NewErrSyntaxIncorrect()
	}
}
