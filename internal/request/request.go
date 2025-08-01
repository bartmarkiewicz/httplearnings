package request

import (
	"errors"
	"fmt"
	"io"
	"learnhttp/internal/headers"
	"strings"
)

const BUFFER_SIZE = 8

type ParsingStatus int

const (
	initialised ParsingStatus = iota
	requestStateParsingHeaders
	done
)

type Request struct {
	RequestLine   RequestLine
	ParsingStatus ParsingStatus
	Headers       headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE, BUFFER_SIZE)
	readToIndex := 0
	req := &Request{
		ParsingStatus: initialised,
		Headers:       headers.NewHeaders(),
	}
	for req.ParsingStatus != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.ParsingStatus != done {
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.ParsingStatus != done {
		parsedBytes, err := r.parseSingleLine(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += parsedBytes
		if parsedBytes == 0 {
			// wait for more bytes
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingleLine(data []byte) (int, error) {
	switch r.ParsingStatus {
	case initialised:
		requestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesParsed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.ParsingStatus = requestStateParsingHeaders
		return bytesParsed, nil
	case requestStateParsingHeaders:
		bytesParsed, isDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if isDone {
			r.ParsingStatus = done
		}
		return bytesParsed, nil
	case done:
		return 0, fmt.Errorf("is done")
	}
	return 0, nil
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
	return &requestLine, len([]byte(line)) + 2, nil
}

func isAllCapitalLetters(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
