package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type PingTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.Memory
}

func (s *PingTestSuite) Context() context.Context {
	return s.ctx
}

func (s *PingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.NewMemory()
}

func TestPingTestSuite(t *testing.T) {
	suite.Run(t, new(PingTestSuite))
}

func (s *PingTestSuite) TestPing() {
	result := s.db.Exec(test.PingValue)
	s.Equal(result.Type, protocol.SimpleString)
	s.Equal(result.Value, "PONG")
}
