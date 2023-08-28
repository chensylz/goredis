package main

import (
	"github.com/chensylz/goredis/config"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server"
	"github.com/chensylz/goredis/internal/server/handler"
	"github.com/chensylz/goredis/internal/storage/databse"
)

func main() {
	logger.SetupLog()
	conf := config.Setup("redis.conf")
	s := server.New(*conf,
		handler.NewServer(
			protocol.NewRESP(),
			databse.New(16),
		),
	)
	s.Run()
}
