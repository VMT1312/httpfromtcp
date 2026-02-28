package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	Code response.StatusCode
	Msg  string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h HandlerError) Write(w io.Writer) {
	_ = response.WriteStatusLine(w, h.Code)
	headers := response.GetDefaultHeaders(len(h.Msg))
	_ = response.WriteHeaders(w, headers)
	w.Write([]byte(h.Msg))
}
