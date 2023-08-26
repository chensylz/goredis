package stringcmd

import (
	"context"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
)

func (s *Cmd) Delete(ctx context.Context, key string) *protocol.ProtoValue {
	value := s.db.Delete(ctx, key)
	if value == nil {
		return response.NilBulk
	}
	return response.NewBulkString(value.(string))
}
