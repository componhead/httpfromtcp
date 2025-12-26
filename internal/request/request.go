package request

import (
	"errors"
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
	line, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(string(line), "\r\n")
	for idx, r := range rows {
		// TODO: temporarily only requestline
		if idx > 0 {
			continue
		}
		requestLine, err := parseRequestLine(r)
		if err != nil {
			return nil, err
		}
		request := Request{
			RequestLine: requestLine,
		}
		return &request, nil
	}
	return nil, errors.New("Error on request from reader")
}

func parseRequestLine(r string) (RequestLine, error) {
	parts := strings.Split(r, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("Malformed request line: invalid number of parts")
	}
	method := parts[0]
	runes := []rune(method)
	for _, v := range runes {
		if v < 65 || v > 90 {
			return RequestLine{}, errors.New("Http method unknown")
		}
	}
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, fmt.Errorf("Http version number wrong: %s\n", parts[2])
	}
	httpVersionNumber := (strings.Split(parts[2], "/"))[1]
	requestLine := RequestLine{
		HttpVersion:   httpVersionNumber,
		RequestTarget: parts[1],
		Method:        method,
	}
	return requestLine, nil
}
