package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	req := string(b)
	request_line, err := parseRequestLine(req)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: request_line,
	}, nil
}

func parseRequestLine(request string) (RequestLine, error) {
	lines := strings.Split(request, "\r\n")
	request_line := lines[0]

	parts := strings.Split(request_line, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version_part := parts[2]

	for i := 0; i < len(method); i++ {
		if method[i] < 'A' || method[i] > 'Z' {
			return RequestLine{}, fmt.Errorf("invalid method")
		}
	}

	version := strings.Split(version_part, "/")[1]
	if version != "1.1" {
		return RequestLine{}, fmt.Errorf("invalid http version")
	}

	return RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}
