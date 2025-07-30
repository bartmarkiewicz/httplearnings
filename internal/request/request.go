package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const BUFFER_SIZE = 8

type ParsingStatus int

const (
	initialised ParsingStatus = iota
	done
)

type Request struct {
	RequestLine   RequestLine
	ParsingStatus ParsingStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, BUFFER_SIZE, BUFFER_SIZE)

	readUpToIndex := 0

	request := Request{ParsingStatus: initialised}

	for request.ParsingStatus != done {
		if cap(buffer) == len(buffer) {
			temp := buffer
			buffer = make([]byte, len(buffer)*2, len(buffer)*2)
			copy(buffer, temp)
		}

		readBytes, err := reader.Read(buffer[readUpToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.ParsingStatus = done
				break
			}
			return nil, err
		}
		readUpToIndex += readBytes
		parsedBytes, err := request.parse(buffer[:readUpToIndex])
		readUpToIndex = readUpToIndex - parsedBytes
		copy(buffer, buffer[parsedBytes:])
		if err != nil {
			return nil, err
		}
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParsingStatus == initialised {
		parsedRequestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		} else if bytesParsed == 0 {
			return 0, nil
		}
		r.RequestLine = *parsedRequestLine
		r.ParsingStatus = done
		return bytesParsed, nil
	} else {
		return 0, fmt.Errorf("request already parsed")
	}
}

func parseRequestLine(lines string) (*RequestLine, int, error) {
	if !strings.Contains(lines, "\r\n") {
		return nil, 0, nil
	}

	line := strings.Split(lines, "\r\n")[0]

	splitWords := strings.Split(line, " ")
	if !isAllCapitalLetters(splitWords[0]) {
		return nil, 0, fmt.Errorf("could not extract method from line: %s", line)
	}

	method := splitWords[0]
	httpVersionString := splitWords[2]
	splitHttpVersion := strings.Split(httpVersionString, "/")

	if len(splitHttpVersion) != 2 || splitHttpVersion[1] != "1.1" || splitHttpVersion[0] != "HTTP" {
		return nil, 0, fmt.Errorf("request must be HTTP/1.1 got %s", httpVersionString)
	}

	requestTarget := splitWords[1]

	requestLine := RequestLine{
		HttpVersion:   splitHttpVersion[1],
		RequestTarget: requestTarget,
		Method:        method,
	}
	return &requestLine, len([]byte(line)), nil
}

func isAllCapitalLetters(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
