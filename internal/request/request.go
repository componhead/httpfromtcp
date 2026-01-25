package request

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/componhead/httpfromtcp/internal/headers"
)

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const bufferSize = 8

type requestState int

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       requestState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		state:       requestStateInitialized,
		RequestLine: RequestLine{},
		Headers:     headers.Headers{},
	}
	for req.state != requestStateDone {
		if readToIndex >= len(buf) && req.state != requestStateDone {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		bytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request: %s\n", err)
				}
				break
			}
			return nil, err
		}
		readToIndex += bytesRead
		if readToIndex > 0 {
			numParsed, err := req.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			if numParsed == 0 && req.state != requestStateDone {
				continue
			} else {
				copy(buf, buf[numParsed:readToIndex])
				readToIndex -= numParsed
			}
		}
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	var totalBytesParsed int
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		totalBytesParsed += n
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	var bytesRead int
	switch r.state {
	case requestStateInitialized:
		{
			line, b, err := parseRequestLine(data)
			if err != nil {
				return b, err
			}
			if b == 0 {
				return 0, nil
			}
			r.RequestLine = *line
			r.state = requestStateParsingHeaders
			bytesRead = b
			break
		}
	case requestStateParsingHeaders:
		{
			i, done, err := r.Headers.Parse(data)
			if err != nil {
				return i, err
			}
			if done == true {
				r.state = requestStateDone
			}
			bytesRead = i
			break
		}
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
	return bytesRead, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	var bytesToRead []byte
	var wholeLine bool = false
	for i := range data {
		if i > 1 && data[i-1] == '\r' && data[i] == '\n' {
			wholeLine = true
			bytesToRead = data[:i-1]
			break
		}
	}
	if !wholeLine {
		return nil, 0, nil
	}
	requestLine, err := requestLineFromString(string(bytesToRead))
	if err != nil {
		return nil, len(bytesToRead), err
	}
	return requestLine, len(bytesToRead) + 2, nil
}

func requestLineFromString(requestLineString string) (*RequestLine, error) {
	parts := strings.Split(requestLineString, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s\n", requestLineString)
	}
	method := parts[0]
	runes := []rune(method)
	for i := range runes {
		if runes[i] < 65 || runes[i] > 90 {
			return nil, errors.New("Http method unknown")
		}
	}
	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("Http version number wrong: %s\n", parts[2])
	}
	httpVersionNumber := (strings.Split(parts[2], "/"))[1]
	requestLine := RequestLine{
		HttpVersion:   httpVersionNumber,
		RequestTarget: parts[1],
		Method:        method,
	}
	return &requestLine, nil
}
