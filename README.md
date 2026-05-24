# ThermalKV

ThermalKV is a small in-memory key-value store written in Go. It supports basic operations, optional TTL (time-to-live) for keys, and a simple write-ahead log (WAL) for persisting `SET` operations.

## Features

- In-memory key-value storage
- Thread-safe access using `sync.RWMutex`
- TTL support for expiring keys
- Background cleaner goroutine that evicts expired keys every second
- WAL persistence for `SET` operations (appends to `data/wal.log`)
- Interactive REPL in `cmd/main.go` (SET/GET/DEL/EXIT)

## Repository Structure

- `go.mod` - Go module declaration
- `cmd/main.go` - Interactive REPL that starts the store, replays WAL, and accepts commands
- `internal/store/store.go` - Store implementation (Set, Get, Delete, SetTTL, StartCleaner, Recover)
- `internal/persistence/persistence.go` - WAL writer and loader
- `data/` - Directory where the WAL file `wal.log` is stored (created at runtime)

## Current Behavior

On startup, the example program in `cmd/main.go`:

1. Creates a new store with `store.NewStore()`
2. Starts the background cleaner with `db.StartCleaner()`
3. Loads WAL lines using `persistence.LoadLogs()` and replays them with `db.Recover(logs)`
4. Enters an interactive REPL reading commands from stdin

The REPL supports the following commands:

- `SET <key> <value>` — store a key/value pair (this writes a `SET` entry to the WAL)
- `GET <key>` — retrieve the value for a key
- `DEL <key>` — remove a key from the store (does not write to WAL)
- `EXIT` — exit the REPL

Note: `SetTTL` exists in the `store` API to set expirations programmatically, but the REPL does not expose a command for it.

## Getting Started

### Prerequisites

- Go 1.26 (the module declares `go 1.26.2`)

### Build

```powershell
cd D:/db/ThermalKV
go build ./cmd
```

### Run

```powershell
go run ./cmd
```

The program will replay any existing `data/wal.log` entries on startup before presenting the prompt.

## Store API

The store exposes the following methods (in `internal/store/store.go`):

- `Set(key string, value string)` — add or update a key-value pair (writes to WAL)
- `Get(key string) (string, bool)` — return the value and existence flag
- `Delete(key string)` — remove a key from the store
- `SetTTL(key string, seconds int)` — set key expiration in seconds (in-memory only)
- `StartCleaner()` — start the background cleaner goroutine
- `Recover(logs []string)` — replay WAL `SET` entries into the store

## Persistence

- WAL file path: `data/wal.log` (entries appended by `persistence.WriteLog`)
- Log entry format: `SET <key> <value>` (single-line entries)
- Helpers:
  - `persistence.LoadLogs()` — reads `data/wal.log` and returns a slice of log lines
  - `store.Recover(logs []string)` — replays `SET` entries into the in-memory store

On startup the example `cmd/main.go` calls `persistence.LoadLogs()` then `db.Recover(logs)` to restore previous `SET` operations.

## Notes

- `Get` performs lazy expiration checks and removes expired keys on access.
- The background cleaner removes expired keys periodically (every second).
- `Delete` does not write to the WAL; only `SET` operations are persisted.
