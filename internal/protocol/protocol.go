package protocol

import (
	"bufio"
)

type Processor interface {
	Decode(reader *bufio.Reader) (interface{}, error)
	Encode(values string) (interface{}, error)
}
