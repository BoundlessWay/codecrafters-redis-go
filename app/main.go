package main

import (
	"fmt"
	"net"
	"os"
)


func main() {
	
	fmt.Println("Redis server starting...")

	// 1. Bind port to listen for incoming connections
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	
	defer l.Close()

	// 2. Accept connection from a client
	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	// 3. Handle client request (for simplicity, we just respond with PONG)
	for {
		buf := make([]byte, 1024)

		// Read data from the socket
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Connection closed or error encountered.")
			break 
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}
