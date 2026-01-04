package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	stateInitialized requestState = iota
	stateDone
)

const bufferSize = 8

type requestState int

type Request struct {
	RequestLine RequestLine
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
	req := Request{
		state:       stateInitialized,
		RequestLine: RequestLine{},
	}
	for req.state != stateDone {
		if readToIndex >= len(buf) && req.state != stateDone {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		bytesRead, err := reader.Read(buf[readToIndex:])
		var ErrorEOF error = nil
		if err != nil {
			if errors.Is(err, io.EOF) {
				ErrorEOF = err
			} else {
				return nil, err
			}
		}
		readToIndex += bytesRead
		if readToIndex > 0 {
			numParsed, err := req.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			if numParsed == 0 && req.state != stateDone {
				continue
			} else {
				copy(buf, buf[numParsed:readToIndex])
				readToIndex -= numParsed
			}
		}
		if ErrorEOF != nil {
			if req.state == stateDone {
				break
			}
			return nil, fmt.Errorf("no valid request line found: %s\n", ErrorEOF.Error())
		}
	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateDone:
		{
			return 0, fmt.Errorf("error: trying to read data in a done state")
		}
	case stateInitialized:
		{
			line, bytesRead, err := parseRequestLine(data)
			if err != nil {
				return bytesRead, err
			}
			if bytesRead == 0 {
				return 0, nil
			}
			r.RequestLine = *line
			r.state = stateDone
			return bytesRead, nil
		}
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	var bytesToRead []byte
	var wholeLine bool = false
	for i := 0; i < len(data); i++ {
		if i > 1 && data[i-1] == '\r' && data[i] == '\n' {
			wholeLine = true
			bytesToRead = data[:i-1]
			break
		}
	}
	if !wholeLine {
		return nil, 0, nil
	}
	requestLine, bytesRead, err := requestLineFromString(bytesToRead)
	if err != nil {
		return nil, bytesRead, err
	}
	return &requestLine, bytesRead + 2, nil
}

func requestLineFromString(bytesToRead []byte) (RequestLine, int, error) {
	parts := strings.Split(string(bytesToRead), " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("poorly formatted request-line: %s\n", string(bytesToRead))
	}
	method := parts[0]
	runes := []rune(method)
	for i := 0; i < len(runes); i++ {
		if runes[i] < 65 || runes[i] > 90 {
			return RequestLine{}, len(method), errors.New("Http method unknown")
		}
	}
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, int(len(parts[2])), fmt.Errorf("Http version number wrong: %s\n", parts[2])
	}
	httpVersionNumber := (strings.Split(parts[2], "/"))[1]
	requestLine := RequestLine{
		HttpVersion:   httpVersionNumber,
		RequestTarget: parts[1],
		Method:        method,
	}
	return requestLine, len(bytesToRead), nil
}
