package main

import (
	"github.com/chensylz/goredis/config"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/server"
	"github.com/chensylz/goredis/internal/server/handler"
)

func main() {
	logger.SetupLog()
	conf := config.Setup("redis.conf")
	s := server.New(*conf, handler.NewEcho())
	s.Run()
}
