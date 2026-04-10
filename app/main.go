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
			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(args[1]), args[1])
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
		prefix, _ := rd.ReadByte()

		if prefix != '$' {
			return nil, fmt.Errorf("expected $")
		}

		size, _ := readInt(rd)
		data := make([]byte, size)
		io.ReadFull(rd, data)
		rd.ReadString('\n') // Bỏ qua CRLF (\r\n)

		args[i] = string(data)
	}
	return args, nil
}

func readInt(rd *bufio.Reader) (int, error) {
	line, err := rd.ReadString('\n')

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(line))
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
