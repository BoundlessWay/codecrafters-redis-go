package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	okResponse       = "+OK\r\n"
	nullBulkResponse = "$-1\r\n"
	msPerSecond      = 1000

	errTTLInvalid = "value is not an integer or out of range"
	errSyntax     = "syntax error"
)

func handleCommand(conn net.Conn, args []string) {
	if len(args) == 0 {
		return
	}

	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "PING":
		writeSimpleString(conn, "PONG")
	case "ECHO":
		if len(args) != 2 {
			writeWrongArgCount(conn, "echo")
			return
		}
		writeBulkString(conn, args[1])
	case "SET":
		handleSet(conn, args)
	case "GET":
		handleGet(conn, args)
	default:
		fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", cmd)
	}
}

func handleSet(conn net.Conn, args []string) {
	if len(args) != 3 && len(args) != 5 {
		writeWrongArgCount(conn, "set")
		return
	}

	key, value := args[1], args[2]
	entry := storeEntry{value: value}

	if len(args) == 5 {
		expireAtMs, errResp := parseExpiryMillis(args[3], args[4])
		if errResp != nil {
			writeError(conn, errResp.Error())
			return
		}
		entry.hasExpiry = true
		entry.expireAtMs = expireAtMs
	}

	storeMu.Lock()
	store[key] = entry
	storeMu.Unlock()
	fmt.Fprint(conn, okResponse)
}

func handleGet(conn net.Conn, args []string) {
	if len(args) != 2 {
		writeWrongArgCount(conn, "get")
		return
	}

	key := args[1]
	storeMu.RLock()
	entry, exists := store[key]
	storeMu.RUnlock()
	if !exists {
		fmt.Fprint(conn, nullBulkResponse)
		return
	}

	nowMs := time.Now().UnixMilli()
	if entry.hasExpiry && nowMs >= entry.expireAtMs {
		// Re-check under write lock to avoid deleting a newly-updated key.
		storeMu.Lock()
		freshEntry, stillExists := store[key]
		if stillExists && freshEntry.hasExpiry && nowMs >= freshEntry.expireAtMs {
			delete(store, key)
			storeMu.Unlock()
			fmt.Fprint(conn, nullBulkResponse)
			return
		}
		storeMu.Unlock()

		if !stillExists {
			fmt.Fprint(conn, nullBulkResponse)
			return
		}
		writeBulkString(conn, freshEntry.value)
		return
	}

	writeBulkString(conn, entry.value)
}

func parseExpiryMillis(option string, rawTTL string) (int64, error) {
	ttlValue, err := strconv.ParseInt(rawTTL, 10, 64)
	if err != nil || ttlValue <= 0 {
		return 0, fmt.Errorf(errTTLInvalid)
	}

	switch strings.ToUpper(option) {
	case "PX":
		return time.Now().UnixMilli() + ttlValue, nil
	case "EX":
		return time.Now().UnixMilli() + (ttlValue * msPerSecond), nil
	default:
		return 0, fmt.Errorf(errSyntax)
	}
}
