package protocol

type Processor interface {
	Decode(input []byte) (interface{}, error)
	Encode(values string) (interface{}, error)
}
