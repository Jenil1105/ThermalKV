package tests

import (
	"os"
	"testing"
	"thermalkv/internal/coldstore"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/store"
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

	wal := walpkg.NewWAL(
		"tests/testdata/wal.log",
		false,
	)

	manager := coldstore.NewManager("tests/testdata/cold.dat")

	db := store.NewStore(
		wal,
		manager,
	)

	db.MaxHotMemory = 1000000

	return db
}
