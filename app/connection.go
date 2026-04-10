package main

import (
	"fmt"
	"net"
	"os"
)

func bind() net.Listener {
	address := net.JoinHostPort(serverHost, serverPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to bind to %s\n", address)
		os.Exit(1)
	}
	return listener
}

func handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleClient(conn)
	}
}
