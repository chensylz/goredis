package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// RESPDataType 表示RESP协议中的数据类型

const (
	SimpleString byte = '+'
	Error        byte = '-'
	Integer      byte = ':'
	BulkString   byte = '$'
	Array        byte = '*'
)

type RESP struct{}

// NewRESP 创建一个RESP协议处理器
func NewRESP() *RESP {
	return &RESP{}
}

func (r *RESP) Decode(reader *bufio.Reader) (val *ProtoValue, err error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch firstByte {
	case '+': // 单行字符串
		return r.decodeSimpleString(reader)
	case '-': // 错误消息
		return r.decodeError(reader)
	case ':': // 整数
		return r.decodeInteger(reader)
	case '$': // Bulk String
		return r.decodeBulkString(reader)
	case '*': // 数组
		return r.decodeArray(reader)
	default:
		return nil, fmt.Errorf("unknown RESP type: %c", firstByte)
	}
}

func (r *RESP) Encode(values *ProtoValue) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := r.encode(buffer, values)
	return buffer.Bytes(), err
}

func (r *RESP) MustEncode(values *ProtoValue) []byte {
	data, err := r.Encode(values)
	if err != nil {
		panic(err)
	}
	return data
}

func (r *RESP) decodeArray(reader *bufio.Reader) (*ProtoValue, error) {
	line, err := r.readLine(reader)
	if err != nil {
		return nil, err
	}
	arrayLength, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}
	if arrayLength == -1 {
		return nil, nil
	}
	if arrayLength < 0 {
		return nil, errors.New("invalid array length")
	}

	array := make([]*ProtoValue, arrayLength)
	for i := 0; i < arrayLength; i++ {
		elemTypeByte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		switch elemTypeByte {
		case '+':
			elem, err := r.decodeSimpleString(reader)
			if err != nil {
				return nil, err
			}
			array[i] = elem
		case '-':
			elem, err := r.decodeError(reader)
			if err != nil {
				return nil, err
			}
			array[i] = elem
		case ':':
			elem, err := r.decodeInteger(reader)
			if err != nil {
				return nil, err
			}
			array[i] = elem
		case '$':
			elem, err := r.decodeBulkString(reader)
			if err != nil {
				return nil, err
			}
			array[i] = elem
		case '*':
			elem, err := r.decodeArray(reader)
			if err != nil {
				return nil, err
			}
			array[i] = elem
		default:
			return nil, fmt.Errorf("unknown element type: %c", elemTypeByte)
		}
	}
	return &ProtoValue{
		Type:  Array,
		Value: array,
	}, nil
}

func (r *RESP) decodeSimpleString(reader *bufio.Reader) (*ProtoValue, error) {
	line, err := r.readLine(reader)
	if err != nil {
		return nil, err
	}
	return &ProtoValue{Type: SimpleString, Value: string(line)}, nil
}

func (r *RESP) decodeError(reader *bufio.Reader) (*ProtoValue, error) {
	line, err := r.readLine(reader)
	if err != nil {
		return nil, err
	}
	return &ProtoValue{
		Type:  Error,
		Value: string(line),
	}, nil
}

func (r *RESP) decodeInteger(reader *bufio.Reader) (*ProtoValue, error) {
	line, err := r.readLine(reader)
	if err != nil {
		return nil, err
	}
	value, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return nil, err
	}
	return &ProtoValue{
		Type:  Integer,
		Value: value,
	}, nil
}

func (r *RESP) decodeBulkString(reader *bufio.Reader) (*ProtoValue, error) {
	line, err := r.readLine(reader)
	if err != nil {
		return nil, err
	}
	stringLength, err := strconv.Atoi(strings.TrimSpace(string(line)))
	if err != nil {
		return nil, err
	}
	if stringLength == -1 {
		return &ProtoValue{
			Type:  BulkString,
			Value: "",
		}, nil
	}
	if stringLength < 0 {
		return nil, errors.New("invalid bulk string length")
	}

	stringContent := make([]byte, stringLength)
	_, err = reader.Read(stringContent)
	if err != nil {
		return nil, err
	}
	_, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return &ProtoValue{
		Type:  BulkString,
		Value: string(stringContent),
	}, nil
}

func (r *RESP) readLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(line), nil
}

func (r *RESP) encode(buffer *bytes.Buffer, data *ProtoValue) error {
	switch data.Type {
	case SimpleString:
		r.encodeSimpleString(buffer, data)
	case Error:
		r.encodeError(buffer, data)
	case Integer:
		r.encodeInteger(buffer, data)
	case Array:
		r.encodeArray(buffer, data)
	default:
		return errors.New("unknown protocol type")
	}
	return nil
}

func (r *RESP) encodeSimpleString(buffer *bytes.Buffer, data *ProtoValue) {
	buffer.WriteString(fmt.Sprintf("+%s\r\n", data.Value))
}

func (r *RESP) encodeError(buffer *bytes.Buffer, data *ProtoValue) {
	buffer.WriteString(fmt.Sprintf("-%s\r\n", data.Value))
}

func (r *RESP) encodeInteger(buffer *bytes.Buffer, data *ProtoValue) {
	buffer.WriteString(fmt.Sprintf(":%v\r\n", data.Value))
}

func (r *RESP) encodeArray(buffer *bytes.Buffer, data *ProtoValue) {
	array := data.Value.([]*ProtoValue)
	buffer.WriteString(fmt.Sprintf("*%d\r\n", len(array)))
	for _, elem := range array {
		r.encode(buffer, elem)
	}
}
