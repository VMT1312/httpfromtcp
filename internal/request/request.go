package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	status      parseState
}

type parseState int

const (
	initialized parseState = iota
	done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	switch r.status {
	case initialized:
		requestLine, byteConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if byteConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.status = done
		return byteConsumed, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := Request{
		RequestLine: RequestLine{},
		status:      initialized,
	}
	for request.status != done {
		if readToIndex == len(buf) {
			oldBuf := buf
			newSize := len(oldBuf) * 2
			buf = make([]byte, newSize)
			_ = copy(buf, oldBuf)
		}
		readBytes, err := reader.Read(buf[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			request.status = done
			break
		}
		readToIndex = readToIndex + readBytes
		parsedBytes, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[parsedBytes:readToIndex])
		readToIndex -= parsedBytes
	}
	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}
	reqestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(reqestLineText)
	if err != nil {
		return nil, idx, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	lines := strings.Split(str, "\r\n")
	request_line := lines[0]

	parts := strings.Split(request_line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version_part := parts[2]

	for i := 0; i < len(method); i++ {
		if method[i] < 'A' || method[i] > 'Z' {
			return nil, fmt.Errorf("invalid method")
		}
	}

	version := strings.Split(version_part, "/")[1]
	if version != "1.1" {
		return nil, fmt.Errorf("invalid http version")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}
