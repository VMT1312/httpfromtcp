package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"strconv"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, h Handler) (*Server, error) {
	addr := strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		return &Server{}, err
	}
	s := &Server{
		listener: listener,
		handler:  h,
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		HandlerError{
			Code: response.StatusBadRequest,
			Msg:  fmt.Sprintf("Bad request: %v", err),
		}.Write(conn)
		return
	}
	s.handler(response.NewWriter(conn), req)
}
