package handler

import (
	"bufio"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/chensylz/goredis/internal/global/serrors"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
	"github.com/chensylz/goredis/internal/storage"
)

type Server struct {
	Connections sync.Map
	Active      atomic.Bool
	Processor   protocol.Processor
	Storage     storage.Storage
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
		data, err := s.Processor.Decode(reader)
		if err != nil {
			if err == io.EOF {
				logger.Infof("client %s close connection", serverConn.Address())
				s.Connections.Delete(serverConn)
			} else {
				logger.Errorf("read message error: %s", err)
				err = serverConn.Write(serrors.ErrProtocol)
				if err != nil {
					logger.Errorf("conn write error: %s", err)
				}
			}
			return
		}
		msg := data.ToCommand()
		logger.Infof("receive message: %s", msg)
		result, err := s.Storage.Exec(serverConn, msg)
		if err != nil {
			logger.Errorf("exec message error: %s", err)
			err = serverConn.Write(serrors.ErrExec)
			if err != nil {
				logger.Errorf("conn write error: %s", err)
			}
			return
		}
		err = serverConn.Write(result)
		if err != nil {
			logger.Errorf("conn write error: %s", err)
			return
		}
	}
}

func (s *Server) Close() {
	s.Active.Store(true)
	s.Connections.Range(func(key, value interface{}) bool {
		conn := key.(*connections.Server)
		_ = conn.Close()
		return true
	})
}
