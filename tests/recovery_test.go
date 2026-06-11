package tests

import (
	"os"
	"testing"
	"thermalkv/internal/persistence/snapshot"
	recoverpkg "thermalkv/internal/recover"
)

func TestRecoverWALRoundTripPreservesLatestState(t *testing.T) {
	source := NewNamedTestEnv(t, "source")
	recovered := NewNamedTestEnv(t, "recovered")

	source.Store.Set("A", "hello world")
	source.Store.Set("B", "persist me")
	source.Store.SetTTL("B", 3600)
	source.Store.Set("C", "old value")
	source.Store.Set("C", "new value again")
	source.Store.Delete("A")

	recoverpkg.RecoverWAL(recovered.Store, source.Dir)

	if _, exists := recovered.Store.Get("A"); exists {
		t.Fatal("expected deleted key A to stay deleted after WAL recovery")
	}

	value, exists := recovered.Store.Get("B")
	if !exists {
		t.Fatal("expected key B to exist after WAL recovery")
	}
	if value != "persist me" {
		t.Fatalf("expected key B value %q, got %q", "persist me", value)
	}

	if recovered.Store.Data["B"].Expiry == 0 {
		t.Fatal("expected recovered key B to retain its expiry")
	}

	value, exists = recovered.Store.Get("C")
	if !exists {
		t.Fatal("expected key C to exist after WAL recovery")
	}
	if value != "new value again" {
		t.Fatalf("expected latest value for key C, got %q", value)
	}

	expectedMemory := int64(len("persist me") + len("new value again"))
	if recovered.Store.CurrentMemoryUsage != expectedMemory {
		t.Fatalf("expected recovered memory usage %d, got %d", expectedMemory, recovered.Store.CurrentMemoryUsage)
	}
}

func TestRecoverSnapshotRoundTripUsesConfiguredPath(t *testing.T) {
	source := NewNamedTestEnv(t, "snapshot_source")
	recovered := NewNamedTestEnv(t, "snapshot_recovered")

	source.Store.Set("name", "ThermalKV")
	source.Store.Set("lang", "Go")
	source.Store.SetTTL("lang", 3600)

	snapshotData := source.Store.ExportData()
	if err := snapshot.SaveSnapshot(source.Paths.SnapshotPath, snapshotData); err != nil {
		t.Fatalf("save snapshot failed: %v", err)
	}

	recoverpkg.RecoverSnapshot(recovered.Store, source.Paths.SnapshotPath)

	value, exists := recovered.Store.Get("name")
	if !exists || value != "ThermalKV" {
		t.Fatalf("expected snapshot to restore name=ThermalKV, got exists=%v value=%q", exists, value)
	}

	value, exists = recovered.Store.Get("lang")
	if !exists || value != "Go" {
		t.Fatalf("expected snapshot to restore lang=Go, got exists=%v value=%q", exists, value)
	}

	if recovered.Store.Data["lang"].Expiry == 0 {
		t.Fatal("expected snapshot recovery to preserve expiry metadata")
	}
}

func TestRecoverColdIndexUsesProvidedPath(t *testing.T) {
	env := NewTestEnv(t)

	content := []byte(
		"keep|value|4102444800\n" +
			"drop|value|1\n" +
			"keep2|other|4102444800\n" +
			"DEL|keep2\n",
	)

	if err := os.WriteFile(env.Paths.ColdPath, content, 0644); err != nil {
		t.Fatalf("failed to seed cold file: %v", err)
	}

	if err := recoverpkg.RecoverColdIndex(env.Manager, env.Paths.ColdPath); err != nil {
		t.Fatalf("recover cold index failed: %v", err)
	}

	if !env.Manager.HaveIndex("keep") {
		t.Fatal("expected unexpired key to be indexed")
	}

	if env.Manager.HaveIndex("drop") {
		t.Fatal("expected expired key to be skipped during recovery")
	}

	if env.Manager.HaveIndex("keep2") {
		t.Fatal("expected DEL marker to remove recovered cold index entry")
	}
}
