package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type GetSetTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.Memory
}

func (s *GetSetTestSuite) Context() context.Context {
	return s.ctx
}

func (s *GetSetTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.NewMemory()
}

func TestGetSetTestSuite(t *testing.T) {
	suite.Run(t, new(GetSetTestSuite))
}

func (s *GetSetTestSuite) TestGetSet() {
	result := s.db.Exec(test.SetValue)
	s.NotEqual(result.Type, protocol.Error)
	result = s.db.Exec(test.GetValue)
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "value")
}
