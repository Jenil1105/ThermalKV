package tests

import (
	"fmt"
	"testing"
)

func BenchmarkSetWithoutCooling(b *testing.B) {
	env := NewTestEnv(b)
	db := env.Store
	db.MaxHotMemory = 1000000000 // 1GB - high memory limit so no cooling triggers

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := "value-data"
		db.Set(key, value)
	}
}

func BenchmarkSetWithCooling(b *testing.B) {
	env := NewTestEnv(b)
	db := env.Store
	db.MaxHotMemory = 1000000000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := "value-data"
		db.Set(key, value)
		_ = db.CoolKey(key)
	}
}

func BenchmarkGetWithoutCooling(b *testing.B) {
	env := NewTestEnv(b)
	db := env.Store
	db.MaxHotMemory = 1000000000

	// Pre-populate keys
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := "value-data"
		db.Set(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, _ = db.Get(key)
	}
}

func BenchmarkGetWithCooling(b *testing.B) {
	env := NewTestEnv(b)
	db := env.Store
	db.MaxHotMemory = 1000000000

	// Pre-populate and cool all keys
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := "value-data"
		db.Set(key, value)
		_ = db.CoolKey(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, _ = db.Get(key)
	}
}
