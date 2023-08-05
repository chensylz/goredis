package storage

import "github.com/chensylz/goredis/internal/protocol"

type Storage interface {
	Exec(args [][]byte) *protocol.ProtoValue
}
