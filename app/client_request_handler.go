package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		prefix, err := reader.ReadByte()
		if err != nil {
			return
		}

		if prefix != '*' {
			fmt.Fprintf(conn, "-ERR Protocol error: expected '*', got '%c'\r\n", prefix)
			reader.ReadString('\n')
			continue
		}

		args, err := parseArrayContent(reader)
		if err != nil {
			fmt.Fprintf(conn, "-ERR Protocol error: %s\r\n", err.Error())
			continue
		}

		handleCommand(conn, args)
	}
}

func handleCommand(conn net.Conn, args []string) {
	if len(args) == 0 {
		return
	}

	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	case "ECHO":
		if len(args) != 2 {
			fmt.Fprintf(conn, "-ERR wrong number of arguments for 'echo' command\r\n")
			return
		}
		fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(args[1]), args[1])
	case "SET":
		if len(args) != 3 && len(args) != 5 {
			fmt.Fprintf(conn, "-ERR wrong number of arguments for 'set' command\r\n")
			return
		}

		key, value := args[1], args[2]

		var expireAt int64
		if len(args) == 5 {
			option := strings.ToUpper(args[3])
			ttlValue, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil || ttlValue <= 0 {
				fmt.Fprintf(conn, "-ERR value is not an integer or out of range\r\n")
				return
			}

			var ttlMs int64
			switch option {
			case "PX":
				ttlMs = ttlValue
			case "EX":
				ttlMs = ttlValue * 1000
			default:
				fmt.Fprintf(conn, "-ERR syntax error\r\n")
				return
			}
			expireAt = time.Now().UnixMilli() + ttlMs
		}

		storeMu.Lock()
		store[key] = value
		if expireAt > 0 {
			expiryAtMs[key] = expireAt
		} else {
			delete(expiryAtMs, key)
		}
		storeMu.Unlock()
		fmt.Fprint(conn, "+OK\r\n")
	case "GET":
		if len(args) != 2 {
			fmt.Fprintf(conn, "-ERR wrong number of arguments for 'get' command\r\n")
			return
		}
		key := args[1]
		storeMu.RLock()
		value, exists := store[key]
		expireAt, hasExpiry := expiryAtMs[key]
		storeMu.RUnlock()
		if !exists {
			fmt.Fprint(conn, "$-1\r\n")
			return
		}

		if hasExpiry && time.Now().UnixMilli() >= expireAt {
			storeMu.Lock()
			delete(store, key)
			delete(expiryAtMs, key)
			storeMu.Unlock()
			fmt.Fprint(conn, "$-1\r\n")
			return
		}

		fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
	default:
		fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", cmd)
	}
}

