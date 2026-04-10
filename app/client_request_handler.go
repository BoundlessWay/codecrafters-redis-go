package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	protocolErrPattern = "-ERR Protocol error: %s\r\n"
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
			fmt.Fprintf(conn, protocolErrPattern, err.Error())
			continue
		}

		handleCommand(conn, args)
	}
}

