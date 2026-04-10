package main

import (
	"fmt"
	"sync"
)

var (
	serverHost = "0.0.0.0"
	serverPort = "6379"
	store      = make(map[string]string)
	storeMu    sync.RWMutex
)

func main() {
	fmt.Println("Redis server starting...")

	listener := bind()
	defer listener.Close()

	handleConnections(listener)
}
