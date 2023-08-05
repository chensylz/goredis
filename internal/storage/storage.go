package storage

import "github.com/chensylz/goredis/internal/server/connections"

type Storage interface {
	Exec(server *connections.Server, args [][]byte) ([]byte, error)
}
