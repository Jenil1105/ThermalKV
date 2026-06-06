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
