package commands

import (
	"context"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
)

type StringCmd interface {
	Get(ctx context.Context, key string) *protocol.ProtoValue
	Set(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	GetSet(ctx context.Context, key string, value interface{}) *protocol.ProtoValue
	Delete(ctx context.Context, key string) *protocol.ProtoValue
}

type ExpireCmd interface {
	Expire(ctx context.Context, key string, expiredAt int64) *protocol.ProtoValue
}

type CommonCmd interface {
	Ping(ctx context.Context) *protocol.ProtoValue
	Echo(ctx context.Context, value string) *protocol.ProtoValue
	Select(ctx context.Context, db string, conn *connections.Server) *protocol.ProtoValue
}
