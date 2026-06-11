package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"thermalkv/internal/persistence/walpkg"
)

func TestWALRotateKeepsConfiguredPath(t *testing.T) {
	dir := filepath.Join("tests", "testdata", sanitizeTestName(t.Name()))
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	walPath := filepath.Join(dir, "custom.log")
	wal := walpkg.NewWAL(walPath, false)
	if wal == nil {
		t.Fatal("failed to create WAL")
	}
	defer wal.Close()

	if err := wal.Write("SET", "alpha", "one"); err != nil {
		t.Fatalf("initial write failed: %v", err)
	}

	rotatedPath, err := wal.Rotate()
	if err != nil {
		t.Fatalf("rotate failed: %v", err)
	}

	if filepath.Dir(rotatedPath) != dir {
		t.Fatalf("expected rotated WAL in %q, got %q", dir, rotatedPath)
	}
	if !strings.HasPrefix(filepath.Base(rotatedPath), "custom_") || filepath.Ext(rotatedPath) != ".log" {
		t.Fatalf("unexpected rotated WAL name: %q", rotatedPath)
	}

	if err := wal.Write("SET", "beta", "two"); err != nil {
		t.Fatalf("write after rotate failed: %v", err)
	}

	activeBytes, err := os.ReadFile(walPath)
	if err != nil {
		t.Fatalf("failed to read active WAL: %v", err)
	}
	if !strings.Contains(string(activeBytes), "SET beta two") {
		t.Fatalf("expected active WAL to contain post-rotation write, got %q", string(activeBytes))
	}

	rotatedBytes, err := os.ReadFile(rotatedPath)
	if err != nil {
		t.Fatalf("failed to read rotated WAL: %v", err)
	}
	if !strings.Contains(string(rotatedBytes), "SET alpha one") {
		t.Fatalf("expected rotated WAL to contain pre-rotation write, got %q", string(rotatedBytes))
	}
}

func TestClearWALUsesProvidedPath(t *testing.T) {
	dir := filepath.Join("tests", "testdata", sanitizeTestName(t.Name()))
	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	walPath := filepath.Join(dir, "wal.log")
	if err := os.WriteFile(walPath, []byte("SET key value\n"), 0644); err != nil {
		t.Fatalf("failed to seed WAL: %v", err)
	}

	walpkg.ClearWAL(walPath)

	info, err := os.Stat(walPath)
	if err != nil {
		t.Fatalf("expected WAL file to exist after clear: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("expected cleared WAL to be empty, size=%d", info.Size())
	}
}
