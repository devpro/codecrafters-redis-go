package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	port := 6379

	fmt.Println("Process starts listening to port", port)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Failed to bind to port", port)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			os.Exit(1)
		}

		//?
		//defer conn.Close()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		smallBuffer := make([]byte, 256)
		readNb, readErr := conn.Read(smallBuffer)
		if readErr != nil {
			if readErr == io.EOF {
				fmt.Println("Received EOF. Stopping loop")
				break
			}
			fmt.Println("Error reading from client: ", readErr.Error())
			continue
		}
		// DEBUG: fmt.Println(smallBuffer[:readNb])

		eolSize := 1
		inputStr := ""
		// issue with telnet on WSL which adds two characters ("\r\n" => 13 10 in byte array)
		if bytes.Contains(smallBuffer, []byte("\r\n")) {
			eolSize = 2
			inputStr = strings.TrimRight(string(smallBuffer[:readNb]), "\r\n")
		} else {
			// another way to get rid of a string in the string
			inputStr = strings.TrimSpace(string(bytes.Replace(smallBuffer[:readNb], []byte("\n"), []byte(""), 1)))
		}
		fmt.Printf("Received \"%s\" (%d bytes)\n", inputStr, readNb-eolSize)

		if strings.Compare(inputStr, "close") == 0 {
			fmt.Println("Received close. Stopping loop")
			break
		}

		_, writeErr := conn.Write([]byte("+PONG\r\n"))
		if writeErr != nil {
			fmt.Println("Error sending to client: ", writeErr.Error())
			continue
		}
	}
}
