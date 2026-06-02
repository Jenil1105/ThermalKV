package store

import (
	"strconv"
	"testing"
	"thermalkv/internal/persistence"
	"thermalkv/internal/thermal"
)

func BenchmarkGet(b *testing.B) {
	println("BENCHMARK STARTED")

	wal := persistence.NewWAL(true)
	defer wal.Close()

	manager := thermal.NewManager()

	s := NewStore(wal, manager)

	s.Set("test", "value")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Get("test")
		}
	})
}

func BenchmarkSet(b *testing.B) {
	wal := persistence.NewWAL(true)
	defer wal.Close()

	manager := thermal.NewManager()

	s := NewStore(wal, manager)

	var counter int64

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter++
			key := "key" + strconv.FormatInt(counter, 10)
			s.Set(key, "value")
		}
	})
}
