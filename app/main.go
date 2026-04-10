package main

import (
	"fmt"
	"net"
	"os"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf);

		if err != nil {
			fmt.Println("Connection closed or error encountered.")
			return
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}


func main() {
	
	fmt.Println("Redis server starting...")

	// 1. Bind port to listen for incoming connections
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	// 2. Accept connections from clients
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		// Event Loop (Netpoller)
		go handleClient(conn)
	}
}
