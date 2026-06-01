package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"thermalkv/internal/persistence"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
)

func main() {

	wal := persistence.NewWAL(true)
	defer wal.Close()

	db := store.NewStore(wal)

	snapshot := persistence.LoadSnapshot()
	db.ImportData(snapshot)

	logs := persistence.LoadLogs()
	db.Recover(logs)

	db.StartCleaner()

	srv, err := server.New(db)

	if err != nil {
		panic(err)
	}

	go srv.Start()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutting down...")
	snapshot = db.ExportData()
	persistence.SaveSnapshot(snapshot)
	srv.Shutdown()
	wal.Close()
	fmt.Println("Shutdown complete")
}
