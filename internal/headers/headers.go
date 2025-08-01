package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	indexOfEndOfLine := bytes.Index(data, []byte(crlf))
	if indexOfEndOfLine == -1 {
		return 0, false, nil
	}
	if indexOfEndOfLine == 0 {
		return 2, true, nil
	}

	headerParts := bytes.SplitN(data[:indexOfEndOfLine], []byte(":"), 2)
	key := string(headerParts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	regexMatch, err := regexp.MatchString(`^[A-Za-z0-9!#$%&'*+\-.^_`+"`"+`|~]+$`, key)
	if err != nil {
		return 0, false, err
	}

	if !regexMatch || len(strings.Trim(key, " ")) <= 1 {
		return 0, false, fmt.Errorf("invalid header: %s", key)
	}

	value := bytes.TrimSpace(headerParts[1])
	key = strings.TrimSpace(key)

	h.Set(key, string(value))
	return indexOfEndOfLine + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if h[strings.ToLower(key)] != "" {
		h[strings.ToLower(key)] += ", " + value
		return
	}

	h[strings.ToLower(key)] = value
}
