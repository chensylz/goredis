package stringcmd

import (
	"context"
	"strconv"

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

func (s *Cmd) Incr(ctx context.Context, key string) *protocol.ProtoValue {
	value := s.db.Get(ctx, key)
	if value == nil {
		s.db.Set(ctx, key, "1")
		return response.One
	}
	i, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil {
		return response.NewErr("don't support incr value")
	}
	s.db.Set(ctx, key, strconv.FormatInt(i+1, 10))
	return response.NewInter(value.(int64) + 1)
}
