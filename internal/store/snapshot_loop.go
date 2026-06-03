package store

import (
	"fmt"
	"os"
	"thermalkv/internal/persistence"
	"time"
)

func (s *Store) StartSnapshotLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			rotatedWal, err := s.WAL.Rotate()

			if err != nil {
				continue
			}

			snapshot := s.ExportData()
			err = persistence.SaveSnapshot(snapshot)

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
