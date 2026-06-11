package tests

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {

	db := NewTestStore(t)

	db.Set("name", "jenil")

	value, exists := db.Get("name")

	if !exists {
		t.Fatal("key should exist")
	}

	if value != "jenil" {
		t.Fatalf(
			"expected jenil, got %s",
			value,
		)
	}
}

func TestDelete(t *testing.T) {

	db := NewTestStore(t)

	db.Set("name", "jenil")

	db.Delete("name")

	_, exists := db.Get("name")

	if exists {
		t.Fatal("key should not exist")
	}
}

func TestTTL(t *testing.T) {

	db := NewTestStore(t)

	db.Set("session", "abc")

	db.SetTTL("session", 1)

	time.Sleep(2 * time.Second)

	_, exists := db.Get("session")

	if exists {
		t.Fatal("key should have expired")
	}
}

func TestOverwriteKey(t *testing.T) {

	db := NewTestStore(t)

	db.Set("name", "jenil")

	db.Set("name", "iitp")

	value, exists := db.Get("name")

	if !exists {
		t.Fatal("key should exist")
	}

	if value != "iitp" {
		t.Fatalf(
			"expected iitp, got %s",
			value,
		)
	}
}

func TestCoolAndLazyRestore(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	db.Set("name", "kept-cold")

	if err := db.CoolKey("name"); err != nil {
		t.Fatalf("cool key failed: %v", err)
	}

	if db.Count() != 0 {
		t.Fatalf("expected no hot keys after cooling, got %d", db.Count())
	}

	if db.CurrentMemoryUsage != 0 {
		t.Fatalf("expected hot memory to drop to 0, got %d", db.CurrentMemoryUsage)
	}

	if !env.Manager.HaveIndex("name") {
		t.Fatal("expected cooled key to remain indexed in cold storage")
	}

	value, exists := db.Get("name")
	if !exists {
		t.Fatal("expected cooled key to be restored on get")
	}

	if value != "kept-cold" {
		t.Fatalf("expected restored value %q, got %q", "kept-cold", value)
	}

	if env.Manager.HaveIndex("name") {
		t.Fatal("expected cold index entry to be removed after restore")
	}

	if db.CurrentMemoryUsage != int64(len("kept-cold")) {
		t.Fatalf("expected hot memory %d after restore, got %d", len("kept-cold"), db.CurrentMemoryUsage)
	}
}

func TestDeleteRemovesColdIndexEntry(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	db.Set("session", "persisted")
	if err := db.CoolKey("session"); err != nil {
		t.Fatalf("cool key failed: %v", err)
	}

	db.Delete("session")

	if env.Manager.HaveIndex("session") {
		t.Fatal("expected cold index entry to be removed after delete")
	}

	if _, exists := db.Get("session"); exists {
		t.Fatal("expected deleted cold key to stay missing")
	}
}

func TestExpiredColdKeyIsDiscardedOnRead(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	db.Set("token", "soon-gone")
	db.SetTTL("token", 1)

	if err := db.CoolKey("token"); err != nil {
		t.Fatalf("cool key failed: %v", err)
	}

	time.Sleep(2 * time.Second)

	if value, exists := db.Get("token"); exists || value != "" {
		t.Fatalf("expected expired cold key to be missing, got exists=%v value=%q", exists, value)
	}

	if env.Manager.HaveIndex("token") {
		t.Fatal("expected expired cold key index to be removed")
	}
}
