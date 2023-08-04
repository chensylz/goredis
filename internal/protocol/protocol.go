package protocol

type Processor interface {
	Parse(input string) (interface{}, error)
	Execute(command interface{}) (interface{}, error)
}
