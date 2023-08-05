package connections

import "net"

type Server struct {
	Conn net.Conn
}

func NewServer(conn net.Conn) *Server {
	return &Server{Conn: conn}
}

func (s *Server) Close() error {
	return s.Conn.Close()
}
