package main

import (
	"thermalkv/internal/persistence"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
)

func main() {

	db := store.NewStore()

	snapshot := persistence.LoadSnapshot()
	db.ImportData(snapshot)

	logs := persistence.LoadLogs()
	db.Recover(logs)

	db.StartCleaner()

	server.Start(db)
}
