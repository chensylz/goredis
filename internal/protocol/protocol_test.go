package protocol_test

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	"github.com/chensylz/goredis/internal/protocol"
	"github.com/stretchr/testify/suite"
)

type ProtocolTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ProtocolTestSuite) Context() context.Context {
	return s.ctx
}

func (s *ProtocolTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

func TestProtocolTestSuite(t *testing.T) {
	suite.Run(t, new(ProtocolTestSuite))
}

func (s *ProtocolTestSuite) TestRESP() {
	respData := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$3\r\nvalue\r\n")
	expectedDecoded := protocol.ProtoValue{
		Type: protocol.Array,
		Value: []protocol.ProtoValue{
			{Type: protocol.BulkString, Value: "SET"},
			{Type: protocol.BulkString, Value: "key"},
			{Type: protocol.BulkString, Value: "value"},
		},
	}
	respProtocol := protocol.NewRESP()
	decoded, err := respProtocol.Decode(bufio.NewReader(bytes.NewReader(respData)))
	s.NoError(err)
	s.Equal(expectedDecoded.Value.([]protocol.ProtoValue)[0].Value, decoded.Value.([]*protocol.ProtoValue)[0].Value)
}
