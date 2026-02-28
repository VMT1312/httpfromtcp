package server

import (
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	Code response.StatusCode
	Msg  string
}

func (h HandlerError) Write(w io.Writer) {
	_ = response.WriteStatusLine(w, h.Code)
	headers := response.GetDefaultHeaders(len(h.Msg))
	_ = response.WriteHeaders(w, headers)
	w.Write([]byte(h.Msg))
}
