# Project Report - Redis Clone (Go)

## 1) Project Overview

This project implements a Redis-like server in Go for Codecrafters.
It handles TCP connections, parses RESP commands, executes core commands, and
returns RESP responses.

## 2) Architecture

- `app/main.go`: server bootstrap and shared store initialization
- `app/connection.go`: bind/listen/accept loop
- `app/client_request_handler.go`: per-client request read loop
- `app/parser.go`: RESP parser
- `app/commands.go`: command dispatch + business logic
- `app/resp_writer.go`: RESP response helpers

### Runtime Flow

1. Server binds on `0.0.0.0:6379`.
2. Each incoming client runs in its own goroutine.
3. Requests are parsed from RESP wire format.
4. Commands execute against in-memory store.
5. Responses are returned in RESP format.

## 3) Implemented Features

- `PING` -> `+PONG\r\n`
- `ECHO <message>`
- `SET <key> <value>`
- `GET <key>`
- `SET <key> <value> PX <milliseconds>`
- `SET <key> <value> EX <seconds>`
- Missing key returns null bulk string (`$-1\r\n`)
- Expired keys are cleaned lazily on `GET`

## 4) Technical Decisions

- Concurrency model: one goroutine per connection
- Shared state protection: `sync.RWMutex`
- Command and parser separation for maintainability
- Case-insensitive option handling for command arguments

## 5) Testing

Test suite is under `tests/` with two modes:

- Raw wire integration tests:
  - `tests/integration_test.go`
  - `tests/harness_redis.go`
- CLI behavior tests:
  - `tests/cli_mode_test.go`

Covered scenarios:

- Success path: `PING`, `ECHO`, `SET`, `GET`
- Error path: wrong arg count, bad TTL, bad option
- Expiry path: immediate read vs post-expiry null response
- Protocol path: malformed RESP input handling

## 6) Current Scope

Implemented for Codecrafters stage progression with focus on:

- correctness
- clarity
- incremental development

Advanced Redis features (replication, persistence, transactions, pub/sub, auth)
are not yet included.
