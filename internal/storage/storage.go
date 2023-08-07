package storage

import (
	"time"

	"github.com/chensylz/goredis/internal/protocol"
)

type Func string

const (
	GET Func = "GET"
	SET Func = "SET"
)

type Entity struct {
	CreatedAt int64
	Size      int64
	Hit       int64
	ExpiredAt int64
	Value     interface{}
}

func NewEntity(value []byte) *Entity {
	return &Entity{
		CreatedAt: time.Now().Unix(),
		Size:      int64(len(value)),
		Hit:       0,
		ExpiredAt: 0,
		Value:     value,
	}
}

type Storage interface {
	Exec(args [][]byte) *protocol.ProtoValue
}
