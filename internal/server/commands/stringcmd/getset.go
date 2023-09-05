package stringcmd

import (
	"context"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Cmd struct {
	db storage.DB
}

func New(db storage.DB) *Cmd {
	return &Cmd{db: db}
}

func (s *Cmd) Get(ctx context.Context, key string) *protocol.ProtoValue {
	value := s.db.Get(ctx, key)
	if value == nil {
		return response.NilBulk
	}
	return response.NewBulkString(value.(string))
}

func (s *Cmd) Set(ctx context.Context, key string, value interface{}) *protocol.ProtoValue {
	strValue := value.(string)
	s.db.Set(ctx, key, strValue)
	return response.Ok
}

func (s *Cmd) GetSet(ctx context.Context, key string, value interface{}) *protocol.ProtoValue {
	str := s.db.Get(ctx, key)
	if str == nil {
		s.db.Set(ctx, key, value)
		return response.NewBulkString(key)
	}
	return response.NewBulkString(str.(string))
}

func (s *Cmd) GetDel(ctx context.Context, key string) *protocol.ProtoValue {
	str := s.db.Get(ctx, key)
	if str != nil {
		s.db.Delete(ctx, key)
	}
	return response.NewBulkString(str.(string))
}
