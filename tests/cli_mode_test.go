package tests

import (
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRedisCLIFlow(t *testing.T) {
	server := exec.Command("go", "run", "./app")
	server.Dir = ".."
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	t.Cleanup(func() {
		if server.Process != nil {
			_ = server.Process.Kill()
			_, _ = server.Process.Wait()
		}
	})

	waitForServerPort(t, "127.0.0.1:6379", 5*time.Second)

	unique := time.Now().Format("150405.000000")
	key := "cli_key_" + unique
	expPxKey := "cli_px_" + unique
	expExKey := "cli_ex_" + unique
	missingKey := "cli_missing_" + unique

	assertCLIEqual(t, "PONG", "PING")
	assertCLIEqual(t, "hello", "ECHO", "hello")
	assertCLIContains(t, "wrong number of arguments", "ECHO")

	assertCLIEqual(t, "OK", "SET", key, "bar")
	assertCLIEqual(t, "bar", "GET", key)
	assertCLINilLike(t, "GET", missingKey)

	assertCLIContains(t, "wrong number of arguments", "SET", "only_key")
	assertCLIContains(t, "syntax error", "SET", "a", "b", "XX", "100")
	assertCLIContains(t, "out of range", "SET", "a", "b", "PX", "0")
	assertCLIContains(t, "wrong number of arguments", "GET")

	assertCLIEqual(t, "OK", "SET", expPxKey, "v", "PX", "120")
	assertCLIEqual(t, "v", "GET", expPxKey)
	time.Sleep(220 * time.Millisecond)
	assertCLINilLike(t, "GET", expPxKey)

	assertCLIEqual(t, "OK", "SET", expExKey, "v2", "EX", "1")
	assertCLIEqual(t, "v2", "GET", expExKey)
	time.Sleep(1200 * time.Millisecond)
	assertCLINilLike(t, "GET", expExKey)
}

func waitForServerPort(t *testing.T, addr string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("server not reachable at %s after %s", addr, timeout)
}

func runRedisCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-p", "6379"}, args...)
	cmd := exec.Command("redis-cli", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("redis-cli failed for args %v: %v, out=%q", args, err, string(out))
	}
	return strings.TrimRight(string(out), "\r\n")
}

func assertCLIEqual(t *testing.T, want string, args ...string) {
	t.Helper()
	got := runRedisCLI(t, args...)
	if got != want {
		t.Fatalf("redis-cli %v: want %q, got %q", args, want, got)
	}
}

func assertCLIContains(t *testing.T, needle string, args ...string) {
	t.Helper()
	got := runRedisCLI(t, args...)
	if !strings.Contains(strings.ToLower(got), strings.ToLower(needle)) {
		t.Fatalf("redis-cli %v: expected to contain %q, got %q", args, needle, got)
	}
}

func assertCLINilLike(t *testing.T, args ...string) {
	t.Helper()
	got := runRedisCLI(t, args...)
	if got == "" || strings.EqualFold(got, "(nil)") {
		return
	}
	t.Fatalf("redis-cli %v: expected nil-like output, got %q", args, got)
}
