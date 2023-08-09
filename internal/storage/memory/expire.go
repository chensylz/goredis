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

func (m *Memory) expire(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) != 2 {
		return serrors.NewErrSyntaxIncorrect()
	}
	key := args[0].Value.(string)
	expiredAt, err := strconv.ParseInt(args[1].Value.(string), 10, 64)
	if err != nil {
		return serrors.NewErrSyntaxIncorrect()
	}
	if expiredAt < time.Now().Unix() {
		m.Lock()
		delete(m.data, key)
		m.Unlock()
		return serrors.NewOk()
	}
	m.Lock()
	defer m.Unlock()
	entity, ok := m.data[key]
	if !ok {
		return serrors.NewErrKeyNotFound()
	}
	entity.ExpiredAt = expiredAt
	m.data[key] = entity
	m.expireMap.Store(key, expiredAt)
	return serrors.NewOk()
}
