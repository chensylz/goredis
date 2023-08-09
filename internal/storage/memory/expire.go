package memory

import (
	"strconv"
	"time"

	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
)

func (m *Memory) scanExpired() {
	m.expireMap.Range(func(key, value interface{}) bool {
		if value.(int64) < time.Now().Unix() {
			m.RWMutex.Lock()
			delete(m.data, key.(string))
			m.RWMutex.Unlock()
			m.expireMap.Delete(key)
		}
		return true
	})
}

func (m *Memory) expire(args [][]byte) *protocol.ProtoValue {
	if len(args) != 2 {
		return serrors.NewErrSyntaxIncorrect()
	}
	expiredAt, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return serrors.NewErrSyntaxIncorrect()
	}
	if expiredAt < time.Now().Unix() {
		m.Lock()
		delete(m.data, string(args[0]))
		m.Unlock()
		return serrors.NewOk()
	}
	m.Lock()
	defer m.Unlock()
	entity, ok := m.data[string(args[0])]
	if !ok {
		return serrors.NewErrKeyNotFound()
	}
	entity.ExpiredAt = expiredAt
	m.data[string(args[0])] = entity
	m.expireMap.Store(string(args[0]), expiredAt)
	return serrors.NewOk()
}
