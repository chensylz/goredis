package test

import (
	"strconv"

	"github.com/chensylz/goredis/internal/protocol"
)

var (
	SetValue = &protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "SET",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
			{
				Type:  protocol.BulkString,
				Value: "value",
			},
		},
	}

	GetValue = &protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "GET",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
		},
	}

	DeleteValue = &protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "DELETE",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
		},
	}

	PingValue = &protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "PING",
			},
		},
	}
)

func GetExpireKey(expireTime int64) *protocol.ProtoValue {
	expiredAtBytes := []byte(strconv.FormatUint(uint64(expireTime), 10))
	return &protocol.ProtoValue{
		Type: protocol.Array,
		Value: []*protocol.ProtoValue{
			{
				Type:  protocol.BulkString,
				Value: "EXPIRE",
			},
			{
				Type:  protocol.BulkString,
				Value: "key",
			},
			{
				Type:  protocol.BulkString,
				Value: string(expiredAtBytes),
			},
		},
	}
}
