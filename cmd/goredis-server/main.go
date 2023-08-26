package main

import (
	"github.com/chensylz/goredis/config"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server"
	"github.com/chensylz/goredis/internal/server/commands/commoncmd"
	"github.com/chensylz/goredis/internal/server/commands/expirecmd"
	"github.com/chensylz/goredis/internal/server/commands/stringcmd"
	"github.com/chensylz/goredis/internal/server/handler"
	"github.com/chensylz/goredis/internal/storage/memory"
)

func main() {
	logger.SetupLog()
	conf := config.Setup("redis.conf")
	db := memory.New()
	s := server.New(*conf,
		handler.NewServer(
			protocol.NewRESP(),
			db,
			commoncmd.New(db),
			stringcmd.New(db),
			expirecmd.New(db),
		),
	)
	s.Run()
}
