package commoncmd

import (
	"context"
	"fmt"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/storage"
)

type Cmd struct {
	db storage.DB
}

func (c *Cmd) Ping(ctx context.Context) *protocol.ProtoValue {
	return response.NewSimpleString("PONG")
}

func (c *Cmd) Echo(ctx context.Context, value string) *protocol.ProtoValue {
	return response.NewBulkString(value)
}

func (c *Cmd) Info(ctx context.Context, value string) *protocol.ProtoValue {
	return response.NewBulkString(fmt.Sprintf(`# Server
redis_version:7.0.12
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:e9d5a718c4613a7a
redis_mode:standalone
os:Linux 3.10.0-1160.el7.x86_64 x86_64
arch_bits:64
monotonic_clock:POSIX clock_gettime
multiplexing_api:epoll
atomicvar_api:c11-builtin
gcc_version:12.2.0
process_id:1
process_supervised:no
run_id:a28e5a4cfe3c1673b254c0fa804b3f54b199096b
tcp_port:6379
server_time_usec:1694574178009581
uptime_in_seconds:4025
uptime_in_days:0
hz:10
configured_hz:10
lru_clock:75361
executable:/data/redis-server
config_file:
io_threads_active:0
    `))
}

func New(db storage.DB) *Cmd {
	return &Cmd{db: db}
}
