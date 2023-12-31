package stringcmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/stringcmd"
	"github.com/chensylz/goredis/internal/storage/memory"
)

type GetSetTestSuite struct {
	suite.Suite
	ctx context.Context

	cmd commands.StringCmd
}

func (s *GetSetTestSuite) Context() context.Context {
	return s.ctx
}

func (s *GetSetTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.cmd = stringcmd.New(memory.NewSyncDict(0))
}

func TestGetSetTestSuite(t *testing.T) {
	suite.Run(t, new(GetSetTestSuite))
}

func (s *GetSetTestSuite) TestGetSet() {
	result := s.cmd.Set(s.ctx, "key", "123")
	s.NotEqual(result.Type, protocol.Error)
	result = s.cmd.Get(s.ctx, "key")
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "123")
}

func (s *GetSetTestSuite) TestGetDel() {
	result := s.cmd.Set(s.ctx, "key", "123")
	s.NotEqual(result.Type, protocol.Error)
	result = s.cmd.GetDel(s.ctx, "key")
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "123")
	result = s.cmd.Get(s.ctx, "key")
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "")
}

func (s *GetSetTestSuite) TestGetIncr() {
	result := s.cmd.Incr(s.ctx, "key")
	s.Equal(result.Type, protocol.Integer)
	s.Equal(result.Value, int64(1))
	s.cmd.Set(s.ctx, "key1", "xxx")
	result = s.cmd.Incr(s.ctx, "key1")
	s.Equal(result.Type, protocol.Error)
}
