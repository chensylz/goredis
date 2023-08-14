package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	suite.Suite
	ctx context.Context

	db *memory.Memory
}

func (s *DeleteTestSuite) Context() context.Context {
	return s.ctx
}

func (s *DeleteTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = memory.NewMemory()
}

func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (s *DeleteTestSuite) TestExpire() {
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
	s.db.Exec(&protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "DELETE",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
		},
	})
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
	s.Equal("", value.Value)
}
