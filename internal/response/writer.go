package response

import (
	"errors"
	"httpfromtcp/internal/headers"
	"io"
)

type Writer struct {
	writer      io.Writer
	writerState state
}

type state int

const (
	writeStatusLine state = iota
	writeHeaders
	writeBody
)

func (w *Writer) WriteStatusLine(StatusCode StatusCode) error {
	err := WriteStatusLine(w.writer, StatusCode)
	if err != nil {
		return err
	}
	w.writerState = writeHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writeHeaders {
		return errors.New("haven't written status line")
	}
	err := WriteHeaders(w.writer, headers)
	if err != nil {
		return err
	}
	w.writerState = writeBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writeBody {
		return 0, errors.New("haven't written headers")
	}
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, err
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: writeStatusLine,
	}
}
