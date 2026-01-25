package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       parseState
	Headers     headers.Headers
	Body        []byte
}

type parseState int

const (
	requestStateInitialized parseState = iota
	requestStateDone
	requestStateParsingHeaders
	requestStateParseBody
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, byteConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if byteConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return byteConsumed, nil
	case requestStateParsingHeaders:
		byteConsumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if byteConsumed == 0 {
			return 0, nil
		}
		if done {
			r.state = requestStateParseBody
			return byteConsumed, nil
		}
		return byteConsumed, nil
	case requestStateParseBody:
		var i int
		contentLength := r.Headers.Get("Content-Length")
		if contentLength == "" {
			r.state = requestStateDone
			return len(data), nil
		}
		cl, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("invalid content length")
		}
		for i < len(data) && len(r.Body) < cl {
			r.Body = append(r.Body, data[i])
			i++
		}
		if len(r.Body) > cl {
			return 0, fmt.Errorf("length of parsed body is greater than content length")
		}
		if len(r.Body) == cl {
			r.state = requestStateDone
			return i, nil
		}
		return i, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
	}
	for request.state != requestStateDone {
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
			readToIndex += readBytes
			_, err := request.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			if request.state == requestStateDone {
				return &request, nil
			} else {
				return nil, fmt.Errorf("incomplete request")
			}
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
