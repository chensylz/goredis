package storage

import "github.com/chensylz/goredis/internal/protocol"

type Func string

const (
	GET Func = "GET"
	SET Func = "SET"
)

type Entity struct {
	CreatedAt int64
	Size      int64
	Hit       int64
	Value     interface{}
}

type Storage interface {
	Exec(args [][]byte) *protocol.ProtoValue
}
