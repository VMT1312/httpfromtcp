package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Printf("Error open file: %v", err)
	}
	defer f.Close()

	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading file: %v", err)
		}
		txt := string(b[:n])
		fmt.Printf("read: %s\n", txt)
	}
}
