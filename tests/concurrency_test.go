package tests

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentSetGetDeleteDistinctKeys(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	const workers = 32

	var start sync.WaitGroup
	start.Add(1)

	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		i := i
		wg.Add(1)

		go func() {
			defer wg.Done()
			start.Wait()

			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			db.Set(key, value)

			got, exists := db.Get(key)
			if !exists {
				errCh <- fmt.Errorf("expected %s to exist after set", key)
				return
			}
			if got != value {
				errCh <- fmt.Errorf("expected %s=%q, got %q", key, value, got)
				return
			}

			db.Delete(key)

			if _, exists := db.Get(key); exists {
				errCh <- fmt.Errorf("expected %s to be deleted", key)
				return
			}
		}()
	}

	start.Done()
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}

	if count := db.Count(); count != 0 {
		t.Fatalf("expected store to be empty after concurrent deletes, got count=%d", count)
	}

	if db.CurrentMemoryUsage != 0 {
		t.Fatalf("expected memory usage to return to 0, got %d", db.CurrentMemoryUsage)
	}
}

func TestConcurrentReadAfterCoolingRestoresSingleKeySafely(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	db.Set("shared", "cold-value")
	if err := db.CoolKey("shared"); err != nil {
		t.Fatalf("cool key failed: %v", err)
	}

	const readers = 24

	var start sync.WaitGroup
	start.Add(1)

	var wg sync.WaitGroup
	errCh := make(chan error, readers)

	for i := 0; i < readers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			start.Wait()

			value, exists := db.Get("shared")
			if !exists {
				errCh <- fmt.Errorf("expected cooled key to be restored for concurrent reader")
				return
			}
			if value != "cold-value" {
				errCh <- fmt.Errorf("expected restored value %q, got %q", "cold-value", value)
			}
		}()
	}

	start.Done()
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}

	if env.Manager.HaveIndex("shared") {
		t.Fatal("expected cold index entry to be removed after concurrent restore")
	}

	if count := db.Count(); count != 1 {
		t.Fatalf("expected exactly one hot key after restore, got %d", count)
	}

	if db.CurrentMemoryUsage != int64(len("cold-value")) {
		t.Fatalf("expected memory usage %d after restore, got %d", len("cold-value"), db.CurrentMemoryUsage)
	}
}
