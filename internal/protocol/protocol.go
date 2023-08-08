package protocol

import (
	"bufio"
	"fmt"
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

func (v *ProtoValue) ToCommand() [][]byte {
	value, ok := v.Value.([]*ProtoValue)
	if ok {
		var command [][]byte
		for _, v := range value {
			command = append(command, v.ToCommand()...)
		}
		return command
	}
	return [][]byte{[]byte(fmt.Sprintf("%s", v.Value))}
}
