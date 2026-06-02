# ThermalKV

ThermalKV is a high-performance key-value database written in Go that combines in-memory speed with durable storage mechanisms. It supports concurrent access, TTL-based expiration, write-ahead logging, snapshot persistence, and thermal storage concepts where data can be moved between hot (memory) and cold (disk) storage tiers.

The project was built to explore core database internals such as storage engines, persistence, recovery, caching strategies, concurrency control, and data lifecycle management.

---

## Features

* In-memory key-value storage with thread-safe access.
* Per-key TTL (Time-To-Live) expiration.
* Background expiration cleaner using a min-heap scheduler.
* Write-Ahead Logging (WAL) for durability and crash recovery.
* Snapshot persistence for faster startup and recovery.
* Cold storage support for moving infrequently used data to disk.
* Lazy restoration of cold data back into memory.
* TCP server/client architecture.
* Concurrent read and write operations using Go's synchronization primitives.
* Graceful shutdown with persistence support.

---

## Architecture

```text
                 +----------------+
                 |     Client     |
                 +----------------+
                          |
                          v
                 +----------------+
                 |   TCP Server   |
                 +----------------+
                          |
                          v
                 +----------------+
                 |     Store      |
                 +----------------+
                    |    |    |
          +---------+    |    +---------+
          |              |              |
          v              v              v
   +------------+ +------------+ +------------+
   | TTL Heap   | |    WAL     | | Snapshot   |
   +------------+ +------------+ +------------+
                                        |
                                        v
                               +----------------+
                               | Cold Storage   |
                               |   (cold.dat)   |
                               +----------------+
```

---

## Storage Tiers

### Hot Storage

* Stored entirely in memory.
* Provides the lowest latency access.
* Supports TTL expiration.
* Handles the majority of read and write operations.

### Cold Storage

* Persisted on disk.
* Used for data that has been manually cooled.
* Reduces memory usage.
* Supports lazy loading back into memory when accessed.

### Future Thermal Vision

ThermalKV is designed with a thermal-storage architecture in mind:

* Hot Keys → Frequently accessed keys kept in memory.
* Warm Keys → Data stored on disk but indexed for faster retrieval.
* Cold Keys → Rarely accessed data persisted on disk.
* Automatic tier migration based on access frequency and inactivity.

---

## Key Components

### Server

`cmd/server/main.go`

* Loads persisted data.
* Replays WAL entries.
* Starts background services.
* Handles TCP connections.
* Supports graceful shutdown.

### Client

`cmd/client/main.go`

* Interactive command-line client.
* Connects to the server over TCP.
* Sends commands and displays responses.

### Store

`internal/store`

Core database engine responsible for:

* Set/Get/Delete operations.
* TTL management.
* Snapshot creation.
* Recovery.
* Cold storage interaction.

### Persistence

`internal/persistence`

Handles:

* WAL writes.
* WAL recovery.
* Snapshot export.
* Snapshot import.

### Thermal Manager

`internal/thermal`

Responsible for:

* Moving keys into cold storage.
* Maintaining cold storage metadata.
* Reloading cooled keys into memory.

### TTL Scheduler

`internal/ttl`

Implements:

* Min-heap expiration queue.
* Efficient expiration scheduling.
* Cleaner wake-up coordination.

---

## Runtime Files

| File                | Purpose              |
| ------------------- | -------------------- |
| `data/wal.log`      | Write-ahead log      |
| `data/snapshot.dat` | Snapshot persistence |
| `data/cold.dat`     | Cold storage data    |

---

## Supported Commands

| Command               | Description                        |
| --------------------- | ---------------------------------- |
| `SET <key> <value>`   | Store or update a value            |
| `GET <key>`           | Retrieve a value                   |
| `DEL <key>`           | Delete a key                       |
| `TTL <key> <seconds>` | Set expiration                     |
| `COOL <key>`          | Move key into cold storage         |
| `COUNT`               | Number of keys currently in memory |
| `EXISTS <key>`        | Check if a key exists              |
| `KEYS`                | List all keys                      |
| `INFO`                | Store statistics                   |
| `EXIT`                | Close client connection            |

---

## Example Session

```text
SET user1 jenil
OK

GET user1
jenil

TTL user1 30
OK

EXISTS user1
true

COOL user1
OK

GET user1
jenil
```

The final `GET` restores the key from cold storage back into memory.

---

## Durability and Recovery

ThermalKV uses two persistence mechanisms:

### Write-Ahead Log (WAL)

Every mutation operation is appended to:

```text
data/wal.log
```

This includes:

* SET
* DEL
* EXPIRE

The WAL ensures operations can be replayed after crashes or unexpected shutdowns.

### Snapshots

Periodic snapshots capture the full in-memory state and store it in:

```text
data/snapshot.dat
```

Snapshots significantly reduce recovery time by avoiding full WAL replay.

---

## Recovery Process

On startup the server:

1. Opens or creates the WAL file.
2. Loads the latest snapshot.
3. Imports snapshot data into memory.
4. Replays WAL entries created after the snapshot.
5. Discards expired records.
6. Starts the TTL cleaner.

This guarantees data consistency after restarts.

---

## Technical Highlights

* Concurrent key-value store built in Go.
* RWMutex-based synchronization.
* Min-heap TTL scheduler.
* Lazy expiration strategy.
* Write-Ahead Logging (WAL).
* Snapshot-based persistence.
* Cold storage migration.
* Lazy restoration of cooled data.
* TCP client/server architecture.
* Crash recovery through WAL replay.

---

## Project Structure

```text
ThermalKV/
│
├── cmd/
│   ├── server/
│   └── client/
│
├── internal/
│   ├── model/
│   ├── persistence/
│   ├── store/
│   ├── thermal/
│   └── ttl/
│
├── data/
│   ├── wal.log
│   ├── snapshot.dat
│   └── cold.dat
│
├── go.mod
└── README.md
```

---

## Build and Run

### Prerequisites

* Go (1.26.2)

### Start Server

```bash
go run ./cmd/server
```

### Start Client

```bash
go run ./cmd/client
```

### Sample Commands

```text
SET name ThermalKV
GET name

TTL name 60

COUNT

INFO

EXIT
```

---

## Current Status

Implemented:

* In-memory storage
* TTL expiration
* WAL persistence
* Snapshot recovery
* Cold storage
* TCP server/client
* Concurrent operations

Planned:

* Automatic hot-to-cold migration
* Warm storage layer
* Access-frequency tracking
* Background compaction
* Disk indexing
* Replication support
* Metrics and observability

---

## Why ThermalKV?

ThermalKV is a learning-oriented database project that explores how modern storage systems manage durability, recovery, memory efficiency, and data temperature. It combines concepts commonly found in production databases—such as WALs, snapshots, caching layers, and storage tiering—into a compact Go implementation suitable for experimentation and extension.
