package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"thermalkv/internal/coldstore"
	"thermalkv/internal/persistence/walpkg"
	"thermalkv/internal/store"
)

func main() {
	numOps := flag.Int("n", 5000, "Number of operations to run for each benchmark scenario")
	flag.Parse()

	fmt.Println("==================================================================")
	fmt.Println(" ThermalKV Throughput and Cooling Benchmark Utility")
	fmt.Println("==================================================================")
	fmt.Printf("Running benchmarks with N = %d operations per scenario...\n\n", *numOps)

	tempDir := filepath.Join("data", "benchmark_run")
	
	// Ensure fresh state
	_ = os.RemoveAll(tempDir)
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create temporary benchmark directory: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		// Final cleanup
		_ = os.RemoveAll(tempDir)
	}()

	// Helper to initialize a store instance
	initStore := func() (*store.Store, *walpkg.WAL, *coldstore.Manager, func()) {
		walPath := filepath.Join(tempDir, "wal.log")
		coldPath := filepath.Join(tempDir, "cold.dat")

		wal := walpkg.NewWAL(walPath, false)
		if wal == nil {
			panic("Failed to create WAL")
		}

		manager := coldstore.NewManager(coldPath)
		if manager == nil || manager.IceTray == nil {
			wal.Close()
			panic("Failed to create cold store manager")
		}

		db := store.NewStore(wal, manager)
		db.MaxHotMemory = 1000000000 // 1GB - high so emergency cooling doesn't interfere

		cleanup := func() {
			wal.Close()
			if manager.IceTray != nil && manager.IceTray.File != nil {
				_ = manager.IceTray.File.Close()
			}
			_ = os.RemoveAll(tempDir)
			_ = os.MkdirAll(tempDir, 0755)
		}

		return db, wal, manager, cleanup
	}

	results := make(map[string]struct {
		duration  time.Duration
		opsPerSec float64
		latency   time.Duration
	})

	runAndMeasure := func(name string, scenario func(db *store.Store)) {
		db, _, _, cleanup := initStore()
		defer cleanup()

		fmt.Printf("Running %s... ", name)
		start := time.Now()
		scenario(db)
		duration := time.Since(start)
		fmt.Println("Done")

		opsPerSec := float64(*numOps) / duration.Seconds()
		avgLatency := duration / time.Duration(*numOps)

		results[name] = struct {
			duration  time.Duration
			opsPerSec float64
			latency   time.Duration
		}{duration, opsPerSec, avgLatency}
	}

	// 1. SET Without Cooling
	runAndMeasure("SET (Without Cooling)", func(db *store.Store) {
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "val-" + strconv.Itoa(i)
			db.Set(key, value)
		}
	})

	// 2. SET With Cooling
	runAndMeasure("SET (With Cooling)", func(db *store.Store) {
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "val-" + strconv.Itoa(i)
			db.Set(key, value)
			_ = db.CoolKey(key)
		}
	})

	// 3. GET Without Cooling (Hot Storage)
	runAndMeasure("GET (Without Cooling)", func(db *store.Store) {
		// Populate first
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "val-" + strconv.Itoa(i)
			db.Set(key, value)
		}
		// Measure GET
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			_, _ = db.Get(key)
		}
	})

	// 4. GET With Cooling (Lazy Restore)
	runAndMeasure("GET (With Cooling)", func(db *store.Store) {
		// Populate and cool first
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "val-" + strconv.Itoa(i)
			db.Set(key, value)
			_ = db.CoolKey(key)
		}
		// Measure GET (Cold restore)
		for i := 0; i < *numOps; i++ {
			key := "key-" + strconv.Itoa(i)
			_, _ = db.Get(key)
		}
	})

	// Print beautiful table of results
	fmt.Println("\n==================================================================")
	fmt.Println("                       BENCHMARK RESULTS")
	fmt.Println("==================================================================")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Scenario\tOperations\tElapsed Time\tThroughput\tAvg Latency")
	fmt.Fprintln(w, "--------\t----------\t------------\t----------\t-----------")

	scenarios := []string{
		"SET (Without Cooling)",
		"SET (With Cooling)",
		"GET (Without Cooling)",
		"GET (With Cooling)",
	}

	for _, sc := range scenarios {
		res := results[sc]
		fmt.Fprintf(w, "%s\t%d\t%v\t%.2f ops/sec\t%v\n",
			sc,
			*numOps,
			res.duration.Round(time.Millisecond),
			res.opsPerSec,
			res.latency.Round(time.Nanosecond),
		)
	}
	w.Flush()
	fmt.Println("==================================================================")
}
