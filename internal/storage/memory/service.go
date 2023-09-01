package memory

import (
	"context"
	"sync"
	"time"

	"github.com/chensylz/goredis/internal/storage"
)

type Dict struct {
	index     int
	data      map[string]*storage.Entity
	expireMap sync.Map
	sync.Mutex
}

func (m *Dict) Get(ctx context.Context, key string) interface{} {
	entity, ok := m.data[key]
	if !ok {
		return nil
	}
	m.data[key].Hit++
	return entity.Value
}

func (m *Dict) Set(ctx context.Context, key string, value interface{}) interface{} {
	m.Lock()
	defer m.Unlock()
	entity := storage.NewEntity(value, 0)
	switch value.(type) {
	case string:
		entity.Size = len(value.(string))
	}
	m.data[key] = entity
	return value
}

func (m *Dict) Delete(ctx context.Context, key string) interface{} {
	m.Lock()
	defer m.Unlock()
	entity, ok := m.data[key]
	if !ok {
		return nil
	}
	delete(m.data, key)
	return entity.Value
}

func (m *Dict) SetExpire(ctx context.Context, key string, expiredAt int64) {
	m.expireMap.Store(key, expiredAt)
}

func (m *Dict) RemoveExpire(ctx context.Context, key string) {
	m.expireMap.Delete(key)
}

func (m *Dict) Exists(ctx context.Context, key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *Dict) scanExpired() {
	m.expireMap.Range(func(key, value interface{}) bool {
		if value.(int64) <= time.Now().Unix() {
			m.Lock()
			delete(m.data, key.(string))
			m.Unlock()
			m.expireMap.Delete(key)
		}
		return true
	})
}

func NewSyncDict(index int) *Dict {
	m := &Dict{data: make(map[string]*storage.Entity), index: index}
	go func() {
		m.scanExpired()
		time.Sleep(50 * time.Millisecond)
	}()
	return m
}
