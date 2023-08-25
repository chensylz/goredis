package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	suite.Suite
	ctx context.Context

	db *Memory
}

func (s *DeleteTestSuite) Context() context.Context {
	return s.ctx
}

func (s *DeleteTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = NewMemory()
}

func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (s *DeleteTestSuite) TestExpire() {
	s.db.Exec(test.SetValue)
	s.db.Exec(test.DeleteValue)
	value := s.db.Exec(test.GetValue)
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("", value.Value)
}
