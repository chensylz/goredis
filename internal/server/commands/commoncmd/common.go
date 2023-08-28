package commoncmd

import (
	"context"
	"strconv"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
	"github.com/chensylz/goredis/internal/storage"
)

type Cmd struct {
	db storage.DB
}

func (c *Cmd) Select(ctx context.Context, db string, conn *connections.Server) *protocol.ProtoValue {
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		return response.NewErr("invalid DB index")
	}
	conn.SetDB(uint8(dbNum))
	return response.Ok
}

func (c *Cmd) Ping(ctx context.Context) *protocol.ProtoValue {
	return response.NewSimpleString("PONG")
}

func (c *Cmd) Echo(ctx context.Context, value string) *protocol.ProtoValue {
	return response.NewBulkString(value)
}

func New(db storage.DB) *Cmd {
	return &Cmd{db: db}
}
