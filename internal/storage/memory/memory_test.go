package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
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
	setCommand := [][]byte{
		[]byte("SET"),
		[]byte("key"),
		[]byte("value"),
	}
	result := s.db.Exec(setCommand)
	s.NotEqual(result.Type, protocol.Error)
	getCommand := [][]byte{
		[]byte("GET"),
		[]byte("key"),
	}
	result = s.db.Exec(getCommand)
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, []byte("value"))
}
