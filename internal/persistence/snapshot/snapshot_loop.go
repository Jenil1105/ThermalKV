package snapshot

import (
	"fmt"
	"os"
	"thermalkv/internal/store"
	"time"
)

func StartSnapshotLoop(s *store.Store, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			rotatedWal, err := s.WAL.Rotate()

			if err != nil {
				continue
			}

			snap := s.ExportData()
			err = SaveSnapshot(snap)

			if err != nil {
				fmt.Println("Snapshout error: ", err)
				continue
			}

			err = os.Remove(rotatedWal)

			if err != nil {
				fmt.Println("wal cleanup failed: ", err)
			}

			fmt.Println("Snapshot saved")
		}
	}()
}
