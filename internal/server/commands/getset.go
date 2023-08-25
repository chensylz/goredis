package commands

import (
	"context"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
	"github.com/chensylz/goredis/internal/storage/memory"
)

type StringCmd struct {
	db storage.DB
}

func (s *StringCmd) Get(ctx context.Context, key string) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}

func (s *StringCmd) Set(ctx context.Context, key string, value interface{}) *protocol.ProtoValue {
	strValue := value.(string)
	s.db.Set(ctx, key, strValue)
	return response.Ok
}

func (s *StringCmd) GetSet(ctx context.Context, key string, value interface{}) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}

func (m *memory.DB) get(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) != 1 {
		return response.SyntaxIncorrectErr
	}
	key := args[0].Value.(string)
	v, ok := m.data[key]
	if !ok {
		return response.NilBulk
	}
	v.Hit++
	m.data[key] = v
	return response.NewBulkString(v.Value.(string))
}

func (m *memory.DB) getSet(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) != 2 {
		return response.SyntaxIncorrectErr
	}
	key := args[0].Value.(string)
	value := args[1].Value.(string)
	v, ok := m.data[key]
	if !ok {
		m.data[key] = storage.NewEntity(value, int64(len(value)))
		return response.NewBulkString(key)
	}
	v.Hit++
	m.data[key] = v
	return response.NewBulkString(v.Value.(string))
}
