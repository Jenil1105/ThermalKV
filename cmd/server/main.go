package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"thermalkv/internal/coldstore"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence/snapshot"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/recover"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
	"time"
)

func main() {

	paths := model.Paths{
		WALPath:      "data/wal.log",
		SnapshotPath: "data/snapshot.dat",
		ColdPath:     "data/cold.dat",
	}

	wal := walpkg.NewWAL(paths.WALPath, false)
	wal.StartSyncLoop()
	defer wal.Close()

	manager := coldstore.NewManager(paths.ColdPath)
	db := store.NewStore(wal, manager)

	recover.RecoverSnapshot(db, paths.SnapshotPath)
	recover.RecoverWAL(db, "data")
	err := recover.RecoverColdIndex(manager, paths.ColdPath)

	if err != nil {
		fmt.Println("Cold recovery failed:", err)
	}

	snapshot.StartSnapshotLoop(db, wal, paths.SnapshotPath, 6*time.Minute)
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
	snap := db.ExportData()
	err = snapshot.SaveSnapshot(paths.SnapshotPath, snap)
	srv.Shutdown()
	wal.Close()
	fmt.Println("Shutdown complete")
}
