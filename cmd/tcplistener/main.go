package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
		for s := range getLinesChannel(conn) {
			fmt.Println(s)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	b := make([]byte, 8)
	chnl := make(chan string)
	go func() {
		buff := []string{}
		for {
			n, err := f.Read(b)
			res := string(b[:n])
			if err != nil {
				if len(buff) > 0 {
					chnl <- strings.Join(buff, "")
				}
				close(chnl)
				f.Close()
				break
			}
			if n > 0 {
				s := strings.Split(res, "\n")
				if len(s) > 0 {
					buff = append(buff, s[0])
					for i := 1; i < len(s); i++ {
						chnl <- strings.Join(buff, "")
						buff = []string{s[i]}
					}
				}
			}
		}
	}()
	return chnl
}
