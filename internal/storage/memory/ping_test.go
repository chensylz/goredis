package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type PingTestSuite struct {
	suite.Suite
	ctx context.Context

	db *Memory
}

func (s *PingTestSuite) Context() context.Context {
	return s.ctx
}

func (s *PingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = NewMemory()
}

func TestPingTestSuite(t *testing.T) {
	suite.Run(t, new(PingTestSuite))
}

func (s *PingTestSuite) TestPing() {
	result := s.db.Exec(test.PingValue)
	s.Equal(result.Type, protocol.SimpleString)
	s.Equal(result.Value, "PONG")
}
