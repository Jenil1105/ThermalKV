package coldstore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"thermalkv/internal/model"
	"time"
)

type ColdEntry struct {
	Offset int64
	Expiry int64
}

type Manager struct {
	ColdIndex map[string]ColdEntry
	Mutex     sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		ColdIndex: make(map[string]ColdEntry),
	}
}

func (m *Manager) MoveToCool(key string, item model.Item) error {

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	file, err := os.OpenFile(
		"data/cold.dat",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	offset, err := file.Seek(0, io.SeekEnd)

	if err != nil {
		return err
	}

	record := fmt.Sprintf("%s|%s|%d\n", key, item.Value, item.Expiry)

	_, err = file.WriteString(record)

	if err != nil {
		return err
	}

	m.ColdIndex[key] = ColdEntry{
		Offset: offset,
		Expiry: item.Expiry,
	}
	return nil
}

func (m *Manager) LoadFromCool(key string) (model.Item, bool) {

	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	return m.LoadFromCoolNoLock(key)

}

func (m *Manager) LoadFromCoolNoLock(key string) (model.Item, bool) {

	entry, exists := m.ColdIndex[key]

	if !exists {
		return model.Item{}, false
	}

	file, err := os.Open("data/cold.dat")

	if err != nil {
		return model.Item{}, false
	}

	defer file.Close()

	_, err = file.Seek(entry.Offset, 0)

	if err != nil {
		return model.Item{}, false
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return model.Item{}, false
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, "|")

	if len(parts) < 3 {
		return model.Item{}, false
	}

	expiry, _ := strconv.ParseInt(parts[2], 10, 64)

	return model.Item{
		Value:          parts[1],
		Expiry:         expiry,
		LastAccessUnix: time.Now().Unix(),
		Size:           int64(len(parts[1])),
	}, true

}

func (m *Manager) AppendDelete(
	key string,
) error {

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	file, err := os.OpenFile(
		"data/cold.dat",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	record := fmt.Sprintf("DEL|%s\n", key)

	_, err = file.WriteString(record)

	return err
}

func (m *Manager) Compact() error {

	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	fmt.Println("Compaction Started...")

	tempFile, err := os.Create("data/cold_new.dat")

	if err != nil {
		return err
	}

	defer tempFile.Close()

	newIndex := make(map[string]ColdEntry)

	var offset int64 = 0

	for key := range m.ColdIndex {
		item, exists := m.LoadFromCoolNoLock(key)
		if !exists {
			continue
		}
		record := fmt.Sprintf(
			"%s|%s|%d\n",
			key,
			item.Value,
			item.Expiry,
		)
		_, err := tempFile.WriteString(record)

		if err != nil {
			return err
		}

		newIndex[key] = ColdEntry{
			Offset: int64(offset),
			Expiry: item.Expiry,
		}
		offset += int64(len(record))
	}

	tempFile.Close()

	// err = os.Remove("data/cold.dat")

	// if err != nil {
	// 	return err
	// }

	err = os.Rename("data/cold_new.dat", "data/cold.dat")

	if err != nil {
		return err
	}

	m.ColdIndex = newIndex

	fmt.Println("Compaction completed")

	return nil
}
