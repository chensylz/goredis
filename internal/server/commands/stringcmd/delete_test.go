package stringcmd_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/stringcmd"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	suite.Suite
	ctx context.Context

	cmd commands.StringCmd
}

func (s *DeleteTestSuite) Context() context.Context {
	return s.ctx
}

func (s *DeleteTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.cmd = stringcmd.New(memory.NewSyncDict(0))
}

func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (s *DeleteTestSuite) TestExpire() {
	s.cmd.Set(s.ctx, "key", "13")
	s.cmd.Delete(s.ctx, "key")
	value := s.cmd.Get(s.ctx, "key")
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("", value.Value)
}
