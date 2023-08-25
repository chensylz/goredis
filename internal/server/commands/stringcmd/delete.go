package stringcmd

import (
	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
)

func (s *Cmd) delete(args []*protocol.ProtoValue) *protocol.ProtoValue {
	if len(args) < 1 {
		return response.SyntaxIncorrectErr
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
