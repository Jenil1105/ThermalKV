# ThermalKV

ThermalKV is a small, thread-safe in-memory key-value store written in Go. It supports TTL-based expirations, periodic background cleanup, a write-ahead log (WAL), and periodic snapshotting to speed up recovery.

## Features

- In-memory key-value storage with `sync.RWMutex` protection
- TTL support and a background cleaner that evicts expired keys
- Write-ahead log (WAL) at `data/wal.log` for `SET`, `DEL`, and `EXPIRE` operations
- Periodic snapshotting (`data/snapshot.dat`) after a fixed number of writes
- Interactive REPL in `cmd/main.go` with basic commands

## Repository Structure

- `go.mod` - Go module declaration
- `cmd/main.go` - Interactive REPL and startup flow (loads snapshot, replays WAL)
- `internal/store/store.go` - Store implementation (Set, Get, Delete, SetTTL, StartCleaner, Recover, snapshot export/import)
- `internal/persistence/persistence.go` - WAL and snapshot read/write helpers
- `internal/model/types.go` - data model types for items and snapshots
- `internal/ttl/heap.go` - min-heap used for efficiently tracking expirations
- `data/` - runtime data directory where `wal.log` and `snapshot.dat` are stored

## Current Behavior / Startup Flow

On startup, `cmd/main.go`:

1. Creates a new store with `store.NewStore()` and starts the background cleaner (`db.StartCleaner()`).
2. Loads the on-disk snapshot with `persistence.LoadSnapshot()` and imports it via `db.ImportData(snapshot)`.
3. Loads WAL entries with `persistence.LoadLogs()` and replays them using `db.Recover(logs)`.
4. Enters an interactive REPL reading commands from stdin.

The REPL supports the following commands:

- `SET <key> <value>` — store a key/value pair (appends `SET` to the WAL)
- `GET <key>` — retrieve the value for a key
- `DEL <key>` — remove a key from the store (appends `DEL` to the WAL)
- `TTL <key> <seconds>` — set an expiration (appends `EXPIRE` to the WAL)
- `EXIT` — exit the REPL

## Persistence details

- WAL file path: `data/wal.log` (appended by `persistence.WriteLog`)
- Snapshot file: `data/snapshot.dat` (written by `persistence.SaveSnapshot`)
- Log entries include `SET <key> <value>`, `DEL <key>`, and `EXPIRE <key> <expiry-unix>`.
- The store periodically saves a snapshot after a number of write operations (currently every 5 write operations). After saving a snapshot, the WAL is cleared to keep recovery fast.

## Store API (high level)

- `Set(key string, value string)` — add or update a key-value pair (writes `SET` to WAL)
- `Get(key string) (string, bool)` — return the value and existence flag (performs lazy expiration)
- `Delete(key string)` — remove a key from the store (writes `DEL` to WAL)
- `SetTTL(key string, seconds int)` — set key expiration in seconds (writes `EXPIRE` to WAL)
- `StartCleaner()` — start the background cleaner goroutine
- `Recover(logs []string)` — replay WAL entries into the store
- `ExportData()` / `ImportData()` — snapshot export/import helpers used by persistence

## Getting Started

### Prerequisites

- Go 1.26 (module declares `go 1.26.2`)

### Build

```powershell
cd D:/db/ThermalKV
go build ./cmd
```

### Run

```powershell
go run ./cmd
```

The program will load any existing `data/snapshot.dat` and replay `data/wal.log` entries on startup before presenting the prompt.

## Notes and Implementation Details

- `Get` performs lazy expiration checks and will remove expired keys on access.
- A min-heap (`internal/ttl`) is used to efficiently schedule expirations; `StartCleaner` sleeps until the next expiration.
- Snapshots are simple text files (`key|value|expiry` per line) written to `data/snapshot.dat`.
- After snapshotting, the WAL is truncated (cleared) to reduce recovery time.

If you'd like, I can also:

- Add example commands in the README and a quick script to seed `data/wal.log` for testing.
- Add a `Makefile` or simple `build.sh`/`build.ps1` to simplify build/run steps.
