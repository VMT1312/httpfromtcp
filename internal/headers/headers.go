package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	minUpper       = 'A'
	maxUpper       = 'Z'
	minLower       = 'a'
	maxLower       = 'z'
	minDigit       = '0'
	maxDigit       = '9'
	allowedSymbols = "!#$%&'*+-.^_`|~"
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
	key = strings.TrimSpace(key)
	for i := 0; i < len(key); i++ {
		if !isValidHeaderChar(byte(key[i])) {
			return 0, false, fmt.Errorf("invalid character")
		}
	}
	key = strings.ToLower(key)
	value := headerText[colonIdx+1:]
	value = strings.TrimSpace(value)
	_, ok := h[key]
	if ok {
		h[key] = h[key] + ", " + value
	} else {
		h[key] = value
	}
	return idx + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func isValidHeaderChar(ch byte) bool {
	if ch >= minUpper && ch <= maxUpper {
		return true
	}
	if ch >= minLower && ch <= maxLower {
		return true
	}
	if ch >= minDigit && ch <= maxDigit {
		return true
	}
	if strings.ContainsRune(allowedSymbols, rune(ch)) {
		return true
	}
	return false
}
