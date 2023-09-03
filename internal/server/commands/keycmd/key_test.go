package keycmd_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/keycmd"
	"github.com/chensylz/goredis/internal/server/commands/stringcmd"
	"github.com/chensylz/goredis/internal/storage/memory"
)

type ExpireTestSuite struct {
	suite.Suite
	ctx context.Context

	cmd commands.KeyCmd
	str commands.StringCmd
}

func (s *ExpireTestSuite) Context() context.Context {
	return s.ctx
}

func (s *ExpireTestSuite) SetupSuite() {
	s.ctx = context.Background()
	db := memory.NewSyncDict(0)
	s.cmd = keycmd.New(db)
	s.str = stringcmd.New(db)
}

func TestExpireTestSuite(t *testing.T) {
	suite.Run(t, new(ExpireTestSuite))
}

func (s *ExpireTestSuite) TestExpire() {
	key := "key"
	value := "value"
	expiredAt := time.Now().Unix()
	s.str.Set(s.ctx, key, value)
	s.cmd.Expire(s.ctx, key, expiredAt)
	v := s.str.Get(s.ctx, key)
	s.Equal(protocol.BulkString, v.Type)
	s.Equal("", v.Value)
}

func (s *ExpireTestSuite) TestExists() {
	key := "key"
	value := "value"
	s.str.Set(s.ctx, key, value)
	v := s.cmd.Exists(s.ctx, key)
	s.Equal(protocol.Integer, v.Type)
	s.Equal(1, v.Value)
}
