package protocol

import (
	"bufio"
)

type Processor interface {
	Decode(reader *bufio.Reader) (DataType, error)
	Encode(values string) (interface{}, error)
}

type DataType interface {
	ToCommand() [][]byte
}
