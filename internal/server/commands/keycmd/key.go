package keycmd

import (
	"context"
	"time"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Cmd struct {
	db storage.DB
}

func New(db storage.DB) *Cmd {
	return &Cmd{db: db}
}

func (s *Cmd) Expire(ctx context.Context, key string, expiredAt int64) *protocol.ProtoValue {
	if expiredAt <= time.Now().Unix() {
		s.db.Delete(ctx, key)
		s.db.RemoveExpire(ctx, key)
		return response.One
	}
	s.db.SetExpire(ctx, key, expiredAt)
	return response.One
}

func (s *Cmd) Exists(ctx context.Context, key string) *protocol.ProtoValue {
	resp := response.NewInter(0)
	if s.db.Exists(ctx, key) {
		resp.Value = 1
	}
	return resp
}
