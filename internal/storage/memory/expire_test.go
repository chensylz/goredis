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
	s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "SET",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
			{
				Type:  protocol.BulkString,
				Value: "value",
			},
		},
	})
	expiredAt := time.Now().Unix()
	expiredAtBytes := []byte(strconv.FormatUint(uint64(expiredAt), 10))
	s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "EXPIRE",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
			{
				Type:  protocol.BulkString,
				Value: string(expiredAtBytes),
			},
		},
	})
	time.Sleep(50 * time.Millisecond)
	value := s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "GET",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
		},
	})
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("value", value.Value)
	time.Sleep(100 * time.Millisecond)
	s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "EXPIRE",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
			{
				Type:  protocol.BulkString,
				Value: string(expiredAtBytes),
			},
		},
	})
	value = s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "GET",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
		},
	})
	s.Equal(protocol.BulkString, value.Type)
	s.Equal("", value.Value)
}
