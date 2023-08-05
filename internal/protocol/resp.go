package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// RESPDataType 表示RESP协议中的数据类型
type RESPDataType byte

const (
	SimpleString RESPDataType = '+'
	Error        RESPDataType = '-'
	Integer      RESPDataType = ':'
	BulkString   RESPDataType = '$'
	Array        RESPDataType = '*'
)

// RESPValue 表示一个RESP协议中的值
type RESPValue struct {
	Type  RESPDataType
	Value interface{}
}

type RESP struct{}

// NewRESP 创建一个RESP协议处理器
func NewRESP() *RESP {
	return &RESP{}
}

func (r *RESP) Decode(reader *bufio.Reader) (interface{}, error) {
	return decodeRESP(reader)
}

func (r *RESP) Encode(values string) (interface{}, error) {
	return encodeRESP(RESPValue{Type: BulkString, Value: values})
}

// decodeRESP 解析RESP协议数据
func decodeRESP(reader *bufio.Reader) (RESPValue, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return RESPValue{}, err
	}

	dataType := RESPDataType(firstByte)
	switch dataType {
	case SimpleString, Error, Integer:
		line, _, err := reader.ReadLine()
		if err != nil {
			return RESPValue{}, err
		}
		return RESPValue{Type: dataType, Value: string(line)}, nil
	case BulkString:
		lengthLine, _, err := reader.ReadLine()
		if err != nil {
			return RESPValue{}, err
		}
		length, err := strconv.Atoi(string(lengthLine))
		if err != nil {
			return RESPValue{}, err
		}
		if length == -1 {
			return RESPValue{Type: BulkString, Value: nil}, nil
		}
		bulkData := make([]byte, length)
		_, err = reader.Read(bulkData)
		if err != nil {
			return RESPValue{}, err
		}
		// Read trailing CRLF
		_, _, err = reader.ReadLine()
		if err != nil {
			return RESPValue{}, err
		}
		return RESPValue{Type: BulkString, Value: string(bulkData)}, nil
	case Array:
		countLine, _, err := reader.ReadLine()
		if err != nil {
			return RESPValue{}, err
		}
		count, err := strconv.Atoi(string(countLine))
		if err != nil {
			return RESPValue{}, err
		}
		arrayData := make([]RESPValue, count)
		for i := 0; i < count; i++ {
			arrayData[i], err = decodeRESP(reader)
			if err != nil {
				return RESPValue{}, err
			}
		}
		return RESPValue{Type: Array, Value: arrayData}, nil
	default:
		return RESPValue{}, errors.New("unknown data type")
	}
}

func encodeRESP(value RESPValue) ([]byte, error) {
	var buf bytes.Buffer
	switch value.Type {
	case SimpleString, Error, Integer:
		buf.WriteByte(byte(value.Type))
		buf.WriteString(fmt.Sprintf("%s\r\n", value.Value))
	case BulkString:
		buf.WriteByte(byte(value.Type))
		if value.Value == nil {
			buf.WriteString("-1\r\n")
		} else {
			bulkStr := value.Value.(string)
			buf.WriteString(fmt.Sprintf("%d\r\n%s\r\n", len(bulkStr), bulkStr))
		}
	case Array:
		buf.WriteByte(byte(value.Type))
		arrayData := value.Value.([]RESPValue)
		buf.WriteString(fmt.Sprintf("%d\r\n", len(arrayData)))
		for _, item := range arrayData {
			itemBytes, err := encodeRESP(item)
			if err != nil {
				return nil, err
			}
			buf.Write(itemBytes)
		}
	}

	return buf.Bytes(), nil
}
