package main

import (
	"thermalkv/internal/persistence"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
)

func main() {

	wal := persistence.NewWAL()
	defer wal.Close()

	db := store.NewStore(wal)

	snapshot := persistence.LoadSnapshot()
	db.ImportData(snapshot)

	logs := persistence.LoadLogs()
	db.Recover(logs)

	db.StartCleaner()

	server.Start(db)
}
