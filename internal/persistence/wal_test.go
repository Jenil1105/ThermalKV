package persistence

import (
	"testing"
)

func BenchmarkWALWrite(b *testing.B) {
	wal := NewWAL(true)
	wal.StartSyncLoop()
	defer wal.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wal.Write("SET", "key", "value")
	}
}
