package main

import (
	"fmt"
	"net"
)

func writeSimpleString(conn net.Conn, value string) {
	fmt.Fprintf(conn, "+%s\r\n", value)
}

func writeBulkString(conn net.Conn, value string) {
	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
}

func writeWrongArgCount(conn net.Conn, command string) {
	fmt.Fprintf(conn, "-ERR wrong number of arguments for '%s' command\r\n", command)
}

func writeError(conn net.Conn, message string) {
	fmt.Fprintf(conn, "-ERR %s\r\n", message)
}
