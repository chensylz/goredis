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

func (c *Cmd) get(ctx context.Context, key string) (map[string]interface{}, *protocol.ProtoValue) {
	v := c.db.Get(ctx, key)
	if v == nil {
		return nil, response.NewBulkString("")
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, response.NewErr("value is not an hash")
	}
	return m, nil
}

func (c *Cmd) HGet(ctx context.Context, key string, field string) *protocol.ProtoValue {
	m, err := c.get(ctx, key)
	if err != nil {
		return err
	}
	value, ok := m[field]
	if !ok {
		return response.NewBulkString("")
	}
	return response.NewBulkString(value.(string))
}

func (c *Cmd) HSet(ctx context.Context, key string, field string, value interface{}) *protocol.ProtoValue {
	m, err := c.get(ctx, key)
	if err != nil {
		return err
	}
	m[field] = value
	c.db.Set(ctx, key, m)
	return response.One
}

func (c *Cmd) HGetAll(ctx context.Context, key string) *protocol.ProtoValue {
	m, err := c.get(ctx, key)
	if err != nil {
		return err
	}
	values := make([]*protocol.ProtoValue, 0)
	for _, value := range m {
		str := response.NewBulkString(value.(string))
		values = append(values, str)
	}
	return response.NewArray(values)
}

func (c *Cmd) HDel(ctx context.Context, key string, field string) *protocol.ProtoValue {
	m, err := c.get(ctx, key)
	if err != nil {
		return err
	}
	_, ok := m[field]
	if ok {
		delete(m, field)
	}
	return response.One
}

func (c *Cmd) HExists(ctx context.Context, key string, field string) *protocol.ProtoValue {
	//TODO implement me
	panic("implement me")
}
