package memory_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
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
	s.db.Exec([][]byte{
		[]byte("SET"),
		[]byte("key"),
		[]byte("value"),
	})
	expiredAt := time.Now().Unix()
	expiredAtBytes := []byte(strconv.FormatUint(uint64(expiredAt), 10))
	s.db.Exec([][]byte{
		[]byte("EXPIRE"),
		[]byte("key"),
		expiredAtBytes,
	})
	time.Sleep(50 * time.Millisecond)
	value := s.db.Exec([][]byte{
		[]byte("GET"),
		[]byte("key"),
	})
	s.Equal(protocol.BulkString, value.Type)
	s.Equal([]byte("value"), value.Value)
	time.Sleep(100 * time.Millisecond)
	s.db.Exec([][]byte{
		[]byte("EXPIRE"),
		[]byte("key"),
		expiredAtBytes,
	})
	value = s.db.Exec([][]byte{
		[]byte("GET"),
		[]byte("key"),
	})
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("", value.Value)
}
