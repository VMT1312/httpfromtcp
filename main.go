package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Printf("Error open file: %v", err)
	}

	ch := getLinesChannel(f)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func(f io.ReadCloser) {
		defer f.Close()
		defer close(ch)
		var line string
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Error reading file: %v", err)
				break
			}
			txt := string(b[:n])
			parts := strings.Split(txt, "\n")
			parts_len := len(parts)
			for i := 0; i < parts_len-1; i++ {
				part := parts[i]
				complete_line := line + part
				ch <- complete_line
				line = ""
			}
			line += parts[parts_len-1]
		}
		if line != "" {
			ch <- line
		}
	}(f)

	return ch
}
