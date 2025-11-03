package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		return
	}
	defer file.Close()
	b := make([]byte, 8)
	line := "read: "
	for {
		n, err := file.Read(b)
		if err == io.EOF {
			if n > 0 {
				fmt.Printf("read: %s\n", b[:n])
			}
			break
		} else if err != nil {
			return
		}
		if n > 0 {
			fmt.Printf("read: %s\n", b[:n])
		}
	}
}
