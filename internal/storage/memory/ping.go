package memory

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
)

func (m *Memory) ping() *protocol.ProtoValue {
	return serrors.NewSimpleString("PONG")
}
