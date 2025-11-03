package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./message.txt")
	if err != nil {
		fmt.Println("Error reading message.txt file")
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	b := make([]byte, 8)
	for {
		_, err := io.ReadFull(reader, b)
		if err == io.EOF {
			fmt.Println("End of file reached")
			break
		} else if err != nil {
			fmt.Println(err)
			fmt.Println("Unknown error reading file bytes")
			break
		} else {
			fmt.Printf("read: %s", b)
		}
	}
}
