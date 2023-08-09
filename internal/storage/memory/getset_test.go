package memory_test

import (
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage/memory"
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
	setCommand := &protocol.ProtoValue{
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
	}
	result := s.db.Exec(setCommand)
	s.NotEqual(result.Type, protocol.Error)
	getCommand := &protocol.ProtoValue{
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
	}
	result = s.db.Exec(getCommand)
	s.Equal(result.Type, protocol.BulkString)
	s.Equal(result.Value, "value")
}
