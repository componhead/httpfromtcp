package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error creating udp address")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error establishing a connection")
	}
	defer conn.Close()
	rdr := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := rdr.ReadString(byte('\n'))
		if err != nil {
			fmt.Println("Error reading from stdin")
		}
		conn.Write([]byte(str))
	}

}
