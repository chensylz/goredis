package handler

import (
	"bufio"
	"context"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chensylz/goredis/internal/global/response"
	"github.com/chensylz/goredis/internal/logger"
	"github.com/chensylz/goredis/internal/protocol"
	"github.com/chensylz/goredis/internal/server/connections"
	"github.com/chensylz/goredis/internal/storage"
	"github.com/chensylz/goredis/internal/storage/databse"
)

type Server struct {
	Connections sync.Map
	Active      atomic.Bool
	Processor   protocol.Processor
	DBs         *databse.Database
}

func NewServer(
	process protocol.Processor,
	storage *databse.Database,
) *Server {
	return &Server{
		Processor: process,
		DBs:       storage,
	}
}

func (s *Server) Handle(ctx context.Context, conn net.Conn) {
	if s.Active.Load() {
		_ = conn.Close()
		return
	}
	serverConn := connections.NewServer(conn, s.DBs)
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
			s.handler(ctx, reader, serverConn)
		}
	}
}

func (s *Server) Exec(ctx context.Context, args *protocol.ProtoValue, conn *connections.Server) *protocol.ProtoValue {
	if err := s.validArgs(ctx, args); err != nil {
		return err
	}
	value := args.Value.([]*protocol.ProtoValue)
	cmd := value[0].Value.(string)
	cmd = strings.ToUpper(cmd)
	arg := value[1].Value.(string)
	switch storage.Func(cmd) {
	case storage.Set:
		return conn.StrCmd.Set(ctx, arg, value[2].Value)
	case storage.Get:
		return conn.StrCmd.Get(ctx, arg)
	case storage.Expire:
		expiredAt, err := s.parseArgs(ctx, storage.Func(cmd), value[2])
		if err != nil {
			return err
		}
		return conn.KeyCmd.Expire(ctx, arg, expiredAt.(int64))
	case storage.Delete:
		return conn.StrCmd.Delete(ctx, arg)
	case storage.GetSet:
		return conn.StrCmd.GetSet(ctx, arg, value[2].Value)

	case storage.Select:
		conn.Select(arg)
		return response.Ok

	case storage.Ping:
		return conn.ComCmd.Ping(ctx)
	case storage.Echo:
		return conn.ComCmd.Echo(ctx, arg)
	case storage.Info:
		return conn.ComCmd.Info(ctx, arg)
	default:
		return response.SyntaxIncorrectErr
	}
}

func (s *Server) handler(ctx context.Context, reader *bufio.Reader, serverConn *connections.Server) {
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
	if data == nil {
		return
	}
	logger.Infof("receive message: %s", data)
	value := s.Exec(ctx, data, serverConn)
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

func (s *Server) parseArgs(ctx context.Context, argsFun storage.Func, value interface{}) (interface{}, *protocol.ProtoValue) {
	switch argsFun {
	case storage.Expire:
		expiredAt, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return nil, response.SyntaxIncorrectErr
		}
		return expiredAt, nil
	default:
		return nil, response.SyntaxIncorrectErr
	}
}

func (s *Server) validArgs(ctx context.Context, args *protocol.ProtoValue) *protocol.ProtoValue {
	value, ok := args.Value.([]*protocol.ProtoValue)
	if !ok {
		return response.SyntaxIncorrectErr
	}
	fullErr := response.NewErr(response.SyntaxIncorrect + protocol.ProtoValueArrToString(value))
	if len(value) < 2 {
		return fullErr
	}
	switch storage.Func(strings.ToUpper(value[0].Value.(string))) {
	case storage.Set, storage.Expire, storage.GetSet:
		if len(value) != 3 {
			return fullErr
		}
	case storage.Get, storage.Delete, storage.Incr:
		if len(value) != 2 {
			return fullErr
		}
	case storage.Ping:
		if len(value) != 1 {
			return fullErr
		}
	}
	return nil
}
