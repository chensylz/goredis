package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type MemoryTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.DB
}

func (s *MemoryTestSuite) Context() context.Context {
	return s.ctx
}

func (s *MemoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.New()
}

func TestMemoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}

func (s *MemoryTestSuite) TestMemory() {
	result := s.db.Set(s.ctx, "key", "value")
	s.NotNil(result)
	result = s.db.Get(s.ctx, "key")
	s.Equal("value", result)
}
