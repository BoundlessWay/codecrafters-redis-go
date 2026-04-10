package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func startServerAndConnect(t *testing.T) net.Conn {
	t.Helper()

	cmd := exec.Command("go", "run", "./app")
	cmd.Dir = ".."
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	})

	conn := waitForServer(t, "127.0.0.1:6379", 5*time.Second)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func waitForServer(t *testing.T, addr string, timeout time.Duration) net.Conn {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			return conn
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("server not reachable at %s after %s", addr, timeout)
	return nil
}

func sendAndReadSingleResponse(conn net.Conn, args []string) (string, error) {
	wire := encodeRESPArray(args)
	if _, err := conn.Write([]byte(wire)); err != nil {
		return "", err
	}
	return readSingleRESP(conn)
}

func encodeRESPArray(args []string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return b.String()
}

func readSingleRESP(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	prefix, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	switch prefix {
	case '+', '-', ':':
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return string(prefix) + line, nil
	case '$':
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		header := string(prefix) + line
		var n int
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &n); err != nil {
			return "", err
		}
		if n < 0 {
			return header, nil
		}
		payload := make([]byte, n+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return "", err
		}
		return header + string(payload), nil
	default:
		rest, _ := reader.ReadString('\n')
		var buf bytes.Buffer
		buf.WriteByte(prefix)
		buf.WriteString(rest)
		return buf.String(), nil
	}
}
