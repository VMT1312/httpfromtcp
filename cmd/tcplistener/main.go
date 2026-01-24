package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
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

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading request line: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s", r.RequestLine.HttpVersion)
	}

}
