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

const Version = "v0.1.0"

func main() {

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ThermalKV ", Version)
		return
	}

	fmt.Println("=================================")
	fmt.Println(" ThermalKV", Version)
	fmt.Println(" Hot/Cold Key-Value Store")
	fmt.Println("=================================")
	fmt.Println()
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Port: %d\n", 8080)
	fmt.Printf("Data Directory: data/\n")
	fmt.Println()

	paths := model.Paths{
		WALPath:      "data/wal.log",
		SnapshotPath: "data/snapshot.dat",
		ColdPath:     "data/cold.dat",
	}
	os.MkdirAll("data", 0755)

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
