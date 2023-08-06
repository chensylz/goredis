package storage

import "github.com/chensylz/goredis/internal/protocol"

type Func string

const (
	GET Func = "GET"
	SET Func = "SET"
)

type Storage interface {
	Exec(args [][]byte) *protocol.ProtoValue
}
