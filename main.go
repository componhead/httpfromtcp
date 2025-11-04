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
	for s := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", s)
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
			if err == io.EOF {
				if n > 0 {
					chnl <- strings.Join(buff, "")
				}
			}
			if err != nil {
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
