package response

import (
	"errors"
	"fmt"
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
	writeTrailers
	writeDone
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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	total := 0
	if w.writerState != writeBody {
		return 0, errors.New("haven't written headers")
	}
	n, err := fmt.Fprintf(w.writer, "%x\r\n", len(p))
	if err != nil {
		return 0, err
	}
	total += n
	n, err = w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	total += n
	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	total += n
	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writeBody {
		return 0, errors.New("haven't written headers")
	}
	w.writerState = writeTrailers
	return w.writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writeTrailers {
		return errors.New("haven't written body")
	}
	for k, v := range h {
		body := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(body))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	w.writerState = writeDone
	return nil
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: writeStatusLine,
	}
}
