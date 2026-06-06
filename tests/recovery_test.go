package tests

import (
	"os"
	"testing"

	"thermalkv/internal/persistence"
	"thermalkv/internal/store"
	"thermalkv/internal/thermal"
)

func TestWALRecovery(t *testing.T) {

	testDir := "tests/testdata/recovery"

	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0755)

	wal := persistence.NewWAL(
		testDir+"/wal.log",
		false,
	)

	db := store.NewStore(
		wal,
		thermal.NewManager(),
	)

	db.MaxHotMemory = 1000000

	db.Set("A", "hello")
	db.Set("B", "world")

	logs := persistence.LoadLogsFromDir(testDir)

	t.Log(logs)

	recoveredStore := store.NewStore(
		persistence.NewWAL(
			testDir+"/recovered.log",
			false,
		),
		thermal.NewManager(),
	)

	recoveredStore.MaxHotMemory = 1000000

	recoveredStore.Recover(logs)

	value, exists := recoveredStore.Get("A")

	if !exists {
		t.Fatal("A should exist after recovery")
	}

	if value != "hello" {
		t.Fatalf(
			"expected hello, got %s",
			value,
		)
	}

	value, exists = recoveredStore.Get("B")

	if !exists {
		t.Fatal("B should exist after recovery")
	}

	if value != "world" {
		t.Fatalf(
			"expected world, got %s",
			value,
		)
	}
}
