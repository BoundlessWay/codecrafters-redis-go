package main

import (
	"fmt"
	"sync"
)

type storeEntry struct {
	value      string
	expireAtMs int64
	hasExpiry  bool
}

var (
	serverHost = "0.0.0.0"
	serverPort = "6379"
	store      = make(map[string]storeEntry)
	storeMu    sync.RWMutex
)

func main() {
	fmt.Println("Redis server starting...")

	listener := bind()
	defer listener.Close()

	handleConnections(listener)
}
