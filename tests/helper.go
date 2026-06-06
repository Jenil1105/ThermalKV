package tests

import (
	"os"
	"testing"

	"thermalkv/internal/persistence"
	"thermalkv/internal/store"
	"thermalkv/internal/thermal"
)

func NewTestStore(t *testing.T) *store.Store {

	os.Remove(
		"tests/testdata/wal.log",
	)

	err := os.MkdirAll(
		"tests/testdata",
		0755,
	)

	if err != nil {
		t.Fatal(err)
	}

	wal := persistence.NewWAL(
		"tests/testdata/wal.log",
		false,
	)

	manager := thermal.NewManager()

	db := store.NewStore(
		wal,
		manager,
	)

	db.MaxHotMemory = 1000000

	return db
}
