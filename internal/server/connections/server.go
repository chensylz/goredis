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

func (s *Server) Address() string {
	return s.Conn.RemoteAddr().String()
}

func (s *Server) Write(b []byte) error {
	_, err := s.Conn.Write(b)
	return err
}
