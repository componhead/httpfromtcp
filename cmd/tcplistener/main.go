package main

import (
	"fmt"
	"net"

	"github.com/componhead/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error creating listener")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error creating listener")
		}
		fmt.Println("Connection accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	}
}
