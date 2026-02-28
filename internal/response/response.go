package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		b := []byte("HTTP/1.1 200 OK\r\n")
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	case StatusBadRequest:
		b := []byte("HTTP/1.1 400 Bad Reqeust\r\n")
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	case StatusInternalError:
		b := []byte("HTTP/1.1 500 Internal Server Error")
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	default:
		code := strconv.FormatInt(int64(statusCode), 10)
		b := []byte("HTTP/1.1 " + code + " \r\n")
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = strconv.FormatInt(int64(contentLen), 10)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		res := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(res))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
