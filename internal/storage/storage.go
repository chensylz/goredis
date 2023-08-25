package storage

import (
	"context"
	"time"

	"github.com/chensylz/goredis/internal/protocol"
)

type Func string

const (
	Get    Func = "GET"
	GetSet Func = "GETSET"
	Set    Func = "SET"
	Expire Func = "EXPIRE"
	Delete Func = "DELETE"
	Ping   Func = "PING"
)

type Entity struct {
	CreatedAt int64
	Size      int
	Hit       int64
	Value     interface{}
}

func NewEntity(value interface{}, size int) *Entity {
	return &Entity{
		CreatedAt: time.Now().Unix(),
		Size:      size,
		Hit:       0,
		Value:     value,
	}
}

type DB interface {
	Get(ctx context.Context, key string) interface{}
	Set(ctx context.Context, key string, value interface{}) interface{}
}

type StringCmd interface {
	Get(ctx context.Context, key string) *protocol.ProtoValue
	Set(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	GetSet(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	Delete(ctx context.Context, key string) *protocol.ProtoValue
}
