package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
			fmt.Println("Parse error:", err)
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
		if len(args) > 1 {
			writeBulkString(conn, args[1])
		}
	default:
		fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", cmd)
	}
}


func parseArrayContent(rd *bufio.Reader) ([]string, error) {
	// Read the number of elements in the array
	count, err := readInt(rd)
	if err != nil {
		return nil, err
	}

	if count <= 0 {
        return []string{}, nil
    }

	args := make([]string, count)

	for i := 0; i < count; i++ {
		// Each element is expected to be a Bulk String
		arg, err := readBulkString(rd)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	return args, nil
}

func readBulkString(rd *bufio.Reader) (string, error) {

	prefix, err := rd.ReadByte()
	if err != nil || prefix != '$' {
		return "", fmt.Errorf("expected '$'")
	}

	size, err := readInt(rd)
	if err != nil {
		return "", err
	}

	// Handle Null Bulk String ($-1)
	if size == -1 {
		return "", nil
	}

	// Read content (Binary Safe)
	data := make([]byte, size)
	if _, err := io.ReadFull(rd, data); err != nil {
		return "", err
	}

	// Skip \r\n at the end of each Bulk String
	rd.ReadString('\n')

	return string(data), nil
}

func readInt(rd *bufio.Reader) (int, error) {
	line, err := rd.ReadString('\n')

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(line))
}

func writeBulkString(conn net.Conn, s string) {
	// Return in the format: $<length>\r\n<data>\r\n
	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(s), s)
}


func main() {
	
	fmt.Println("Redis server starting...")

	// Bind port to listen for incoming connections
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	// Accept connections from clients
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleClient(conn)
	}
}
