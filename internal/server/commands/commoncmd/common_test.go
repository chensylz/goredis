package commoncmd_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/commoncmd"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type PingTestSuite struct {
	suite.Suite
	ctx context.Context

	cmd commands.CommonCmd
}

func (s *PingTestSuite) Context() context.Context {
	return s.ctx
}

func (s *PingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.cmd = commoncmd.New(memory.New())
}

func TestPingTestSuite(t *testing.T) {
	suite.Run(t, new(PingTestSuite))
}

func (s *PingTestSuite) TestPing() {
	result := s.cmd.Ping(s.ctx)
	s.Equal(result.Type, protocol.SimpleString)
	s.Equal(result.Value, "PONG")
}
