package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/chensylz/goredis/test"
	"github.com/stretchr/testify/suite"
)

type ExpireTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.Memory
}

func (s *ExpireTestSuite) Context() context.Context {
	return s.ctx
}

func (s *ExpireTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.NewMemory()
}

func TestExpireTestSuite(t *testing.T) {
	suite.Run(t, new(ExpireTestSuite))
}

func (s *ExpireTestSuite) TestExpire() {
	s.db.Exec(test.GetValue)
	expiredAt := time.Now().Unix()
	s.db.Exec(test.GetExpireKey(expiredAt))
	time.Sleep(50 * time.Millisecond)
	value := s.db.Exec(test.GetValue)
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("value", value.Value)
	time.Sleep(100 * time.Millisecond)
	s.db.Exec(test.GetExpireKey(expiredAt))
	value = s.db.Exec(test.GetValue)
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("", value.Value)
}
