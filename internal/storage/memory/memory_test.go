package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type MemoryTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.Memory
}

func (s *MemoryTestSuite) Context() context.Context {
	return s.ctx
}

func (s *MemoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.NewMemory()
}

func TestMemoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}

func (s *MemoryTestSuite) TestMemory() {
	result := s.db.Exec(test.SetValue)
	s.NotEqual(result.Type, protocol.Error)
	result = s.db.Exec(test.GetValue)
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "value")
}
