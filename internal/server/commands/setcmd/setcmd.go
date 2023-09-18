package setcmd

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

func (c *Cmd) HGet(ctx context.Context, key string, field string) *protocol.ProtoValue {
	v := c.db.Get(ctx, key)
	if v == nil {
		return response.NewBulkString("")
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return response.NewErr("value is not an hash")
	}
	value, ok := m[field]
	if !ok {
		return response.NewBulkString("")
	}
	return response.NewBulkString(value.(string))
}

func (c *Cmd) HSet(ctx context.Context, key string, field string, value interface{}) *protocol.ProtoValue {
	var (
		m  map[string]interface{}
		ok bool
	)
	v := c.db.Get(ctx, key)
	if v == nil {
		m = make(map[string]interface{})
	}
	m, ok = v.(map[string]interface{})
	if !ok {
		return response.NewErr("value is not an hash")
	}
	m[field] = value
	c.db.Set(ctx, key, m)
	return response.Ok
}

func (c *Cmd) HGetAll(ctx context.Context, key string) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}

func (c *Cmd) HDel(ctx context.Context, key string, field string) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}

func (c *Cmd) HExists(ctx context.Context, key string, field string) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}
