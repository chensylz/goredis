package commands

import (
	"context"

	"github.com/chensylz/goredis/internal/protocol"
)

type StringCmd interface {
	Get(ctx context.Context, key string) *protocol.ProtoValue
	Set(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	GetSet(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	Delete(ctx context.Context, key string) *protocol.ProtoValue
	GetDel(ctx context.Context, key string) *protocol.ProtoValue
	Incr(ctx context.Context, key string) *protocol.ProtoValue
}

type KeyCmd interface {
	Expire(ctx context.Context, key string, expiredAt int64) *protocol.ProtoValue
	Exists(ctx context.Context, key string) *protocol.ProtoValue
}

type CommonCmd interface {
	Ping(ctx context.Context) *protocol.ProtoValue
	Echo(ctx context.Context, value string) *protocol.ProtoValue
	Info(ctx context.Context, value string) *protocol.ProtoValue
}

type SetCmd interface {
	HGet(ctx context.Context, key string, field string) *protocol.ProtoValue
	HSet(ctx context.Context, key string, field string, value interface{}) *protocol.ProtoValue
	HGetAll(ctx context.Context, key string) *protocol.ProtoValue
	HDel(ctx context.Context, key string, field string) *protocol.ProtoValue
	HExists(ctx context.Context, key string, field string) *protocol.ProtoValue
}
