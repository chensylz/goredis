package protocol_test

import (
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
	respData := []byte("*3\r\n$5\r\nHello\r\n:123\r\n-Error\r\n")
	expectedDecoded := protocol.RESPValue{
		Type: protocol.Array,
		Value: []protocol.RESPValue{
			{Type: protocol.BulkString, Value: "Hello"},
			{Type: protocol.Integer, Value: "123"},
			{Type: protocol.Error, Value: "Error"},
		},
	}
	respProtocol := protocol.NewRESP()
	decoded, err := respProtocol.Decode(respData)
	s.NoError(err)
	s.Equal(expectedDecoded, decoded)
}
