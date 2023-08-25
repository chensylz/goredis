package handler

import (
	"context"
	"net"

	"github.com/chensylz/goredis/internal/protocol"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Exec(ctx context.Context, args *protocol.ProtoValue) *protocol.ProtoValue
	Close()
}
