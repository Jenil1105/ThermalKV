# ThermalKV

ThermalKV is a small in-memory key-value store written in Go. It supports basic operations such as `SET`, `GET`, `DELETE`, and `SETTTL`, with a background cleaner that expires keys automatically.

## Features

- In-memory key-value storage
- Thread-safe access using `sync.RWMutex`
- Time-to-live (TTL) support for expiring keys
- Simple write-ahead log (WAL) persistence for `SET` operations
- Example concurrent usage in `cmd/main.go`

## Repository Structure

- `go.mod` - Go module definition
- `cmd/main.go` - Example application that demonstrates store usage and concurrency
- `internal/store/store.go` - In-memory store implementation with TTL and cleaner
- `internal/store/persistence/persistence.go` - Simple WAL persistence for writes
- `data/` - Storage directory for generated WAL file (`wal.log`)

## Getting Started

### Prerequisites

- Go 1.26 or newer

### Build

```bash
cd d:/db/ThermalKV
go build ./cmd
```

### Run

```bash
go run ./cmd
```

This will start the example program in `cmd/main.go`, which creates 500 keys concurrently, performs reads, applies TTL to one key, and deletes others.

## Usage

The store API currently exposes the following methods:

- `Set(key string, value string)` - store a new key-value pair
- `Get(key string) (string, bool)` - retrieve a value and existence flag
- `Delete(key string)` - remove a key
- `SetTTL(key string, seconds int)` - set expiration for an existing key
- `StartCleaner()` - start a background cleaner goroutine to evict expired keys

## Persistence

The store writes each `SET` operation to `data/wal.log` using a very simple append-only log format.

## Notes

- The current persistence layer only logs `SET` operations and does not replay WAL on startup.
- The cleaner runs every second and removes expired keys from memory.
- The `data/` directory should exist or be created by the user for WAL logging to work correctly.

