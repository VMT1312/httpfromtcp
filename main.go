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
		fmt.Printf("Error: %v", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error establishing a connect: %v", err)
			break
		}

		fmt.Println("A connection has been accepted.")

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Connection has been closed.")
	}

}

func getLinesChannel(conn net.Conn) <-chan string {
	ch := make(chan string)
	go func() {
		defer conn.Close()
		defer close(ch)
		var line string
		for {
			b := make([]byte, 8)
			n, err := conn.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error reading file: %v", err)
				continue
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
	}()

	return ch
}
