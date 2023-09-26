package setcmd_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/keycmd"
	"github.com/chensylz/goredis/internal/server/commands/setcmd"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type SetCmdTestSuite struct {
	suite.Suite
	ctx context.Context

	cmd commands.KeyCmd
	set commands.SetCmd
}

func (s *SetCmdTestSuite) Context() context.Context {
	return s.ctx
}

func (s *SetCmdTestSuite) SetupSuite() {
	s.ctx = context.Background()
	db := memory.NewSyncDict(0)
	s.cmd = keycmd.New(db)
	s.set = setcmd.New(db)
}

func TestSetCmdTestSuite(t *testing.T) {
	suite.Run(t, new(SetCmdTestSuite))
}

func (s *SetCmdTestSuite) TestSet() {
	s.set.HSet(s.ctx, "key", "value", "123")
	result := s.set.HGet(s.ctx, "key", "value")
	s.Equal("123", result.Value)
}
