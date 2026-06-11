package snapshot

import (
	"fmt"
	"os"
	"thermalkv/internal/store"
	"time"
)

type WALRotator interface {
	Rotate() (string, error)
}

func StartSnapshotLoop(s *store.Store, wal WALRotator, snapshotPath string, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			rotatedWal, err := wal.Rotate()

			if err != nil {
				continue
			}

			snap := s.ExportData()
			err = SaveSnapshot(snapshotPath, snap)

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
