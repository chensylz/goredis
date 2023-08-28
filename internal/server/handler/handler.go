package handler

import (
	"context"
	"net"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Exec(ctx context.Context, args *protocol.ProtoValue, conn *connections.Server) *protocol.ProtoValue
	Close()
}
