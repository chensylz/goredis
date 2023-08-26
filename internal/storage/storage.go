package storage

import (
	"context"
	"time"
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
	Delete(ctx context.Context, key string) interface{}
	SetExpire(ctx context.Context, key string, expiredAt int64)
	RemoveExpire(ctx context.Context, key string)
}
