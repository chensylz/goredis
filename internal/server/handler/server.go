package handler

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/chensylz/goredis/internal/global/response"
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

func NewServer(process protocol.Processor, storage storage.Storage) *Server {
	return &Server{
		Processor: process,
		Storage:   storage,
	}
}

func (s *Server) Handle(ctx context.Context, conn net.Conn) {
	if s.Active.Load() {
		_ = conn.Close()
		return
	}
	serverConn := connections.NewServer(conn)
	s.Connections.Store(serverConn, struct{}{})
	logger.Infof("client %s connected", serverConn.Address())
	reader := bufio.NewReader(conn)
	for {
		_, ok := s.Connections.Load(serverConn)
		if !ok {
			return
		}
		select {
		case <-ctx.Done():
			_ = conn.Close()
			return
		default:
			s.handler(reader, serverConn)
		}
	}
}

func (s *Server) handler(reader *bufio.Reader, serverConn *connections.Server) {
	data, err := s.Processor.Decode(reader)
	if err != nil {
		if err == io.EOF {
			logger.Infof("client %s close connection", serverConn.Address())
			s.Connections.Delete(serverConn)
		} else {
			logger.Errorf("read message error: %s", err)
			err = serverConn.Write(s.Processor.MustEncode(response.ProtocolErr))
			if err != nil {
				logger.Errorf("conn write error: %s", err)
			}
		}
		return
	}
	logger.Infof("receive message: %s", data)
	value := s.Storage.Exec(data)
	result, err := s.Processor.Encode(value)
	if err != nil {
		logger.Errorf("encode message error: %s", err)
		err = serverConn.Write(s.Processor.MustEncode(response.UnknownErr))
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

func (s *Server) Close() {
	logger.Info("server closing...")
	s.Active.Store(true)
	s.Connections.Range(func(key, value interface{}) bool {
		conn := key.(*connections.Server)
		_ = conn.Close()
		return true
	})
	logger.Info("server closed")
}
