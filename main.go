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
	defer f.Close()

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
			fmt.Printf("read: %s\n", complete_line)
			line = ""
		}
		line += parts[parts_len-1]
	}
	if line != "" {
		fmt.Printf("read: %s\n", line)
	}
}
