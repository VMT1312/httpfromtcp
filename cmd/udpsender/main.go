package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, error := net.ResolveUDPAddr("udp", "localhost:42069")
	if error != nil {
		fmt.Printf("error establishing network: %v", error)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("error dialing: %v", error)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")

		txt, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("err reading: %v", err)
			break
		}

		_, err = conn.Write([]byte(txt))
		if err != nil {
			log.Printf("err writing: %v", err)
			break
		}
	}
}
