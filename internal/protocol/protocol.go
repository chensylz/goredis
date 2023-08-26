package protocol

import (
	"bufio"
	"fmt"
	"strings"
)

type Processor interface {
	Decode(reader *bufio.Reader) (*ProtoValue, error)
	Encode(values *ProtoValue) ([]byte, error)
	MustEncode(values *ProtoValue) []byte
}

type ProtoValue struct {
	Type  byte
	Value interface{}
}

func ProtoValueArrToString(values []*ProtoValue) string {
	resp := make([]string, len(values))
	for i, v := range values {
		resp[i] = fmt.Sprintf("%v", v.Value)
	}
	return strings.Join(resp, " ")
}
