package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Process starts listening to port 6379")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for {
		if _, err := conn.Read([]byte{}); err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}
