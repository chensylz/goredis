package commoncmd

import (
	"context"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Cmd struct {
	db storage.DB
}

func (c *Cmd) Ping(ctx context.Context) *protocol.ProtoValue {
	return response.NewSimpleString("PONG")
}

func New(db storage.DB) *Cmd {
	return &Cmd{db: db}
}
