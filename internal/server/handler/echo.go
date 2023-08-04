package handler

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/chensylz/goredis/internal/logger"
)

type Echo struct {
	Connections sync.Map
	Active      atomic.Bool
}

func NewEcho() *Echo {
	return &Echo{}
}

type EchoConn struct {
	Conn net.Conn
}

func (e *Echo) Handle(ctx context.Context, conn net.Conn) {
	if e.Active.Load() {
		_ = conn.Close()
		return
	}
	echoConn := &EchoConn{conn}
	e.Connections.Store(echoConn, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("client close connection")
				e.Connections.Delete(echoConn)
			} else {
				logger.Error("read message error", err)
			}
			return
		}
		b := []byte(msg)
		_, _ = conn.Write(b)
	}
}

func (e *Echo) Close() {
	e.Active.Store(true)
	e.Connections.Range(func(key, value interface{}) bool {
		conn := key.(*EchoConn)
		_ = conn.Conn.Close()
		return true
	})
}
