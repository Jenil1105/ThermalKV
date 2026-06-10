package coldstore

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"thermalkv/internal/coldstore/index"
)

func (m *Manager) Compact() error {

	m.ColdStore.Lock()
	defer m.ColdStore.Unlock()
	fmt.Println("Compaction Started...")

	tempFile, err := os.Create("data/cold_new.dat")

	if err != nil {
		return err
	}

	defer tempFile.Close()

	newIndex := make(map[string]index.ColdEntry)

	var new_offset int64 = 0

	m.ColdIndex.Mutex.RLock()

	keys := m.ColdIndex.Keys()

	for _, key := range keys {
		old_offset, exists := m.ColdIndex.MeltColdEntryNoLock(key)
		if !exists {
			continue
		}
		cube := m.IceTray.MeltCube(old_offset)
		_, err := tempFile.WriteString(cube)

		if err != nil {
			return err
		}
		parts := strings.Split(cube, "|")
		if len(parts) < 3 {
			continue
		}

		expiry, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}

		newIndex[key] = index.ColdEntry{
			Offset: int64(new_offset),
			Expiry: expiry,
		}
		new_offset += int64(len(cube))
	}
	m.ColdIndex.Mutex.RUnlock()

	tempFile.Close()

	m.IceTray.ChangeIceTray("data/cold_new.dat")

	m.ColdIndex.Mutex.Lock()
	m.ColdIndex.ReplaceEntries(newIndex)
	m.ColdIndex.Mutex.Unlock()

	fmt.Println("Compaction completed")

	return nil
}
