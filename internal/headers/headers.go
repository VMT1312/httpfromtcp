package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	headerText := string(data[:idx])
	colonIdx := strings.Index(headerText, ":")
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("invalid header")
	}
	if colonIdx > 0 {
		if headerText[colonIdx-1:colonIdx] == " " {
			return 0, false, fmt.Errorf("white space before the colon")
		}
	}
	key := headerText[:colonIdx]
	value := headerText[colonIdx+1:]
	h[strings.TrimSpace(key)] = strings.TrimSpace(value)
	return idx + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
