# Graph Report - F:/Workspace/PROJECT/codecrafters-redis-go  (2026-04-11)

## Corpus Check
- Corpus is ~2,160 words - fits in a single context window. You may not need a graph.

## Summary
- 47 nodes · 65 edges · 8 communities detected
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 6 edges (avg confidence: 0.78)
- Token cost: 0 input · 0 output

## God Nodes (most connected - your core abstractions)
1. `Source file roles (architecture)` - 7 edges
2. `TestRedisCLIFlow()` - 5 edges
3. `readBulkString()` - 4 edges
4. `runRedisCLI()` - 4 edges
5. `Wire integration and CLI tests` - 4 edges
6. `handleCommand()` - 3 edges
7. `handleSet()` - 3 edges
8. `parseArrayContent()` - 3 edges
9. `readInt()` - 3 edges
10. `assertCLIEqual()` - 3 edges

## Surprising Connections (you probably didn't know these)
- `Codecrafters Redis clone (Go)` --semantically_similar_to--> `TCP RESP server overview`  [INFERRED] [semantically similar]
  README.md → docs\PROJECT_REPORT.md

## Hyperedges (group relationships)
- **app/ server components** — main, connection, client_request_handler, parser, commands, resp_writer [EXTRACTED 1.00]
- **tests/ verification modes** — integration_test, harness_redis, cli_mode_test [EXTRACTED 1.00]

## Communities

### Community 0 - "Server core and connections"
Cohesion: 0.22
Nodes (2): storeEntry, Source file roles (architecture)

### Community 1 - "CLI integration tests"
Cohesion: 0.62
Nodes (6): assertCLIContains(), assertCLIEqual(), assertCLINilLike(), runRedisCLI(), TestRedisCLIFlow(), waitForServerPort()

### Community 2 - "Raw wire harness"
Cohesion: 0.53
Nodes (5): encodeRESPArray(), readSingleRESP(), sendAndReadSingleResponse(), startServerAndConnect(), waitForServer()

### Community 3 - "Command dispatch"
Cohesion: 0.7
Nodes (4): handleCommand(), handleGet(), handleSet(), parseExpiryMillis()

### Community 4 - "RESP parsing"
Cohesion: 0.8
Nodes (4): parseArrayContent(), readBulkString(), readCRLF(), readInt()

### Community 5 - "RESP encoding"
Cohesion: 0.4
Nodes (0): 

### Community 6 - "Scenarios and feature docs"
Cohesion: 0.4
Nodes (3): Implemented Redis-like commands, Features not yet implemented, Wire integration and CLI tests

### Community 7 - "README and runtime narrative"
Cohesion: 0.4
Nodes (5): Codecrafters Redis clone (Go), RESP parsing and core commands, Goroutine-per-connection and RWMutex, TCP RESP server overview, Five-step request pipeline

## Knowledge Gaps
- **3 isolated node(s):** `storeEntry`, `Goroutine-per-connection and RWMutex`, `Features not yet implemented`
  These have ≤1 connection - possible missing edges or undocumented components.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Source file roles (architecture)` connect `Server core and connections` to `Command dispatch`, `RESP parsing`, `RESP encoding`, `README and runtime narrative`?**
  _High betweenness centrality (0.688) - this node is a cross-community bridge._
- **Why does `TCP RESP server overview` connect `README and runtime narrative` to `Server core and connections`?**
  _High betweenness centrality (0.516) - this node is a cross-community bridge._
- **Why does `Wire integration and CLI tests` connect `Scenarios and feature docs` to `CLI integration tests`, `Raw wire harness`?**
  _High betweenness centrality (0.515) - this node is a cross-community bridge._
- **What connects `storeEntry`, `Goroutine-per-connection and RWMutex`, `Features not yet implemented` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._