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
	requestBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLines := strings.Split(string(requestBytes), "\r\n")
	//for _, line := range requestLines {
	requestLine, err := parseRequestLine(requestLines[0])

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	request := Request{
		RequestLine: *requestLine,
	}
	return &request, nil
	//}
}

func parseRequestLine(line string) (*RequestLine, error) {
	splitWords := strings.Split(line, " ")
	if !isAllCapitalLetters(splitWords[0]) {
		return nil, fmt.Errorf("could not extract method from line: %s", line)
	}

	method := splitWords[0]
	httpVersionString := splitWords[2]
	splitHttpVersion := strings.Split(httpVersionString, "/")

	if len(splitHttpVersion) != 2 || splitHttpVersion[1] != "1.1" || splitHttpVersion[0] != "HTTP" {
		return nil, fmt.Errorf("request must be HTTP/1.1 got %s", httpVersionString)
	}

	requestTarget := splitWords[1]

	requestLine := RequestLine{
		HttpVersion:   splitHttpVersion[1],
		RequestTarget: requestTarget,
		Method:        method,
	}
	return &requestLine, nil
}

func isAllCapitalLetters(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
