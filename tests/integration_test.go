package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRedisScenarios(t *testing.T) {
	conn := startServerAndConnect(t)

	now := time.Now().UnixNano()
	key1 := fmt.Sprintf("k_%d_1", now)
	key2 := fmt.Sprintf("k_%d_2", now)
	key3 := fmt.Sprintf("k_%d_3", now)
	missingKey := fmt.Sprintf("missing_%d", now)

	testCases := []struct {
		name     string
		args     []string
		wantWire string
	}{
		{"ping", []string{"PING"}, "+PONG\r\n"},
		{"echo success", []string{"ECHO", "hello"}, "$5\r\nhello\r\n"},
		{"echo wrong args", []string{"ECHO"}, "-ERR wrong number of arguments for 'echo' command\r\n"},
		{"set/get success set", []string{"SET", key1, "bar"}, "+OK\r\n"},
		{"set/get success get", []string{"GET", key1}, "$3\r\nbar\r\n"},
		{"get missing key", []string{"GET", missingKey}, "$-1\r\n"},
		{"set wrong args", []string{"SET", "only-key"}, "-ERR wrong number of arguments for 'set' command\r\n"},
		{"set bad option", []string{"SET", "a", "b", "XX", "100"}, "-ERR syntax error\r\n"},
		{"set bad ttl", []string{"SET", "a", "b", "PX", "0"}, "-ERR value is not an integer or out of range\r\n"},
		{"get wrong args", []string{"GET"}, "-ERR wrong number of arguments for 'get' command\r\n"},
		{"unknown command", []string{"NOPE"}, "-ERR unknown command 'NOPE'\r\n"},
		{"set px", []string{"SET", key2, "v2", "PX", "50"}, "+OK\r\n"},
		{"set ex", []string{"SET", key3, "v3", "EX", "1"}, "+OK\r\n"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := sendAndReadSingleResponse(conn, tc.args)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if got != tc.wantWire {
				t.Fatalf("response mismatch\nwant: %q\ngot : %q", tc.wantWire, got)
			}
		})
	}

	t.Run("px expires", func(t *testing.T) {
		time.Sleep(80 * time.Millisecond)
		got, err := sendAndReadSingleResponse(conn, []string{"GET", key2})
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if got != "$-1\r\n" {
			t.Fatalf("expected null bulk string after PX expiry, got %q", got)
		}
	})

	t.Run("ex expires", func(t *testing.T) {
		time.Sleep(1100 * time.Millisecond)
		got, err := sendAndReadSingleResponse(conn, []string{"GET", key3})
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if got != "$-1\r\n" {
			t.Fatalf("expected null bulk string after EX expiry, got %q", got)
		}
	})

	t.Run("protocol wrong first byte", func(t *testing.T) {
		raw := "PING\r\n"
		if _, err := conn.Write([]byte(raw)); err != nil {
			t.Fatalf("failed to write raw protocol: %v", err)
		}
		got, err := readSingleRESP(conn)
		if err != nil {
			t.Fatalf("failed to read response: %v", err)
		}
		if !strings.Contains(got, "Protocol error: expected '*'") {
			t.Fatalf("expected protocol error for missing array prefix, got %q", got)
		}
	})

	t.Run("protocol expected dollar", func(t *testing.T) {
		raw := "*1\r\n+PING\r\n"
		if _, err := conn.Write([]byte(raw)); err != nil {
			t.Fatalf("failed to write raw protocol: %v", err)
		}
		got, err := readSingleRESP(conn)
		if err != nil {
			t.Fatalf("failed to read response: %v", err)
		}
		if !strings.Contains(got, "Protocol error: expected '$'") {
			t.Fatalf("expected protocol error for non-bulk arg, got %q", got)
		}
	})
}

