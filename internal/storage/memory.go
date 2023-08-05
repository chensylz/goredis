package storage

import "github.com/chensylz/goredis/internal/global/serrors"

type Memory struct {
	data map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string]interface{})}
}

func (m *Memory) Exec(commands [][]byte) []byte {
	switch string(commands[0]) {
	case "SET":
		return m.set(commands[1:])
	case "GET":
		return m.get(commands[1:])
	default:
		return serrors.ErrSyntaxIncorrect
	}
}

func (m *Memory) set(args [][]byte) []byte {
	if len(args) != 2 {
		return serrors.ErrSyntaxIncorrect
	}
	m.data[string(args[0])] = args[1]
	return serrors.Ok
}

func (m *Memory) get(args [][]byte) []byte {
	if len(args) != 1 {
		return serrors.ErrSyntaxIncorrect
	}
	if v, ok := m.data[string(args[0])]; ok {
		return v.([]byte)
	}
	return serrors.NilBulk
}
