package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		return &Server{}, err
	}
	s := &Server{
		listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Printf("error accepting connection: %v", err)
			break
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!"
	_, err := io.WriteString(conn, resp)
	if err != nil {
		log.Fatalf("error writing to connection: %v", err)
		return
	}
}
