package memory

import (
	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
)

func (m *DB) ping() *protocol.ProtoValue {
	return response.NewSimpleString("PONG")
}
