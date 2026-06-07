package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"thermalkv/internal/coldstore"
	"thermalkv/internal/persistence/snapshot"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/recover"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
	"time"
)

func main() {

	wal := walpkg.NewWAL("data/wal.log", false)
	wal.StartSyncLoop()
	defer wal.Close()

	manager := coldstore.NewManager()
	err := recover.RecoverColdIndex(manager)

	if err != nil {
		fmt.Println("Cold recovery failed:", err)
	}

	db := store.NewStore(wal, manager)
	snapshot.StartSnapshotLoop(db, 6*time.Minute)

	recover.RecoverSnapshot(db, "data/snapshot.dat")

	recover.RecoverWAL(db, "data")

	db.StartCleaner()
	db.StartCoolingWorker()

	srv, err := server.New(db)

	if err != nil {
		panic(err)
	}

	go srv.Start()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutting down...")
	// snap = db.ExportData()
	// snapshot.SaveSnapshot(snap)
	srv.Shutdown()
	wal.Close()
	fmt.Println("Shutdown complete")
}
