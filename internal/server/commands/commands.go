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
}
