package memory

import (
	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/protocol"
)

func (m *Memory) delete(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) < 1 {
		return serrors.NewErrSyntaxIncorrect()
	}
	var count int
	for _, arg := range args {
		if _, ok := m.data[arg.Value.(string)]; ok {
			count++
		}
		delete(m.data, arg.Value.(string))
	}
	return &protocol.ProtoValue{
		Type:  protocol.Integer,
		Value: count,
	}
}
