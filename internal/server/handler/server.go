package handler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
)

type Server struct {
	Connections sync.Map
	Active      atomic.Bool
	Processor   protocol.Processor
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Handle(conn net.Conn) {
	if s.Active.Load() {
		_ = conn.Close()
		return
	}
	serverConn := connections.NewServer(conn)
	s.Connections.Store(serverConn, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := s.Processor.Decode(reader)
		if err != nil {
			if err == io.EOF {
				logger.Info("client close connection")
				s.Connections.Delete(serverConn)
			} else {
				logger.Error("read message error", err)
			}
			return
		}
		fmt.Println(msg)
	}
}
