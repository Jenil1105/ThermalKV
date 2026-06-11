package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"thermalkv/internal/coldstore"
	"thermalkv/internal/model"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/store"
)

type TestEnv struct {
	Dir     string
	Paths   model.Paths
	WAL     *walpkg.WAL
	Manager *coldstore.Manager
	Store   *store.Store
}

func NewTestStore(t *testing.T) *store.Store {
	return NewTestEnv(t).Store
}

func NewTestEnv(t *testing.T) *TestEnv {
	return newTestEnv(t, "")
}

func NewNamedTestEnv(t *testing.T, suffix string) *TestEnv {
	return newTestEnv(t, suffix)
}

func newTestEnv(t *testing.T, suffix string) *TestEnv {
	t.Helper()

	dirName := sanitizeTestName(t.Name())
	if suffix != "" {
		dirName += "_" + sanitizeTestName(suffix)
	}

	dir := filepath.Join("tests", "testdata", dirName)
	_ = os.RemoveAll(dir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	paths := model.Paths{
		WALPath:      filepath.Join(dir, "wal.log"),
		SnapshotPath: filepath.Join(dir, "snapshot.dat"),
		ColdPath:     filepath.Join(dir, "cold.dat"),
	}

	wal := walpkg.NewWAL(paths.WALPath, false)
	if wal == nil {
		t.Fatal("failed to create WAL")
	}

	manager := coldstore.NewManager(paths.ColdPath)
	if manager == nil || manager.IceTray == nil {
		t.Fatal("failed to create cold store manager")
	}

	db := store.NewStore(wal, manager)
	db.MaxHotMemory = 1000000

	t.Cleanup(func() {
		if wal != nil {
			wal.Close()
		}
		if manager != nil && manager.IceTray != nil && manager.IceTray.File != nil {
			_ = manager.IceTray.File.Close()
		}
		_ = os.RemoveAll(dir)
	})

	return &TestEnv{
		Dir:     dir,
		Paths:   paths,
		WAL:     wal,
		Manager: manager,
		Store:   db,
	}
}

func sanitizeTestName(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		" ", "_",
	)
	name = replacer.Replace(name)
	if name == "" {
		return "test"
	}
	return name
}
