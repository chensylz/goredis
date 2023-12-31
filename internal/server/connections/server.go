package connections

import (
	"net"
	"strconv"

	"github.com/chensylz/goredis/internal/server/commands"
	"github.com/chensylz/goredis/internal/server/commands/commoncmd"
	"github.com/chensylz/goredis/internal/server/commands/keycmd"
	"github.com/chensylz/goredis/internal/server/commands/setcmd"
	"github.com/chensylz/goredis/internal/server/commands/stringcmd"
	"github.com/chensylz/goredis/internal/storage"
	"github.com/chensylz/goredis/internal/storage/databse"
)

type Server struct {
	Conn      net.Conn
	dbs       *databse.Database
	dbIndex   uint8
	currentDB storage.DB

	StrCmd commands.StringCmd
	ComCmd commands.CommonCmd
	KeyCmd commands.KeyCmd
	SetCmd commands.SetCmd
}

func NewServer(conn net.Conn, dbs *databse.Database) *Server {
	cDB := dbs.Get(0)
	return &Server{
		Conn:      conn,
		dbs:       dbs,
		currentDB: cDB,
		dbIndex:   0,
		StrCmd:    stringcmd.New(cDB),
		ComCmd:    commoncmd.New(cDB),
		KeyCmd:    keycmd.New(cDB),
		SetCmd:    setcmd.New(cDB),
	}
}

func (s *Server) Close() error {
	return s.Conn.Close()
}

func (s *Server) Address() string {
	return s.Conn.RemoteAddr().String()
}

func (s *Server) Select(arg string) {
	db, err := strconv.Atoi(arg)
	if err != nil {
		return
	}
	s.dbIndex = uint8(db)
	s.currentDB = s.dbs.Get(uint8(db))

	s.StrCmd = stringcmd.New(s.currentDB)
	s.ComCmd = commoncmd.New(s.currentDB)
	s.KeyCmd = keycmd.New(s.currentDB)
	s.SetCmd = setcmd.New(s.currentDB)
}

func (s *Server) Write(b []byte) error {
	_, err := s.Conn.Write(b)
	return err
}
