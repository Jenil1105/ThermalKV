package thermal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"thermalkv/internal/model"
	"time"
)

type ColdEntry struct {
	Offset int64
	Expiry int64
}

type Manager struct {
	ColdIndex map[string]ColdEntry
}

func NewManager() *Manager {
	return &Manager{
		ColdIndex: make(map[string]ColdEntry),
	}
}

func (m *Manager) MoveToCool(key string, item model.Item) error {
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

func (m *Manager) RecoverColdIndex() error {
	file, err := os.Open(
		"data/cold.dat",
	)
	if err != nil {
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var offset int64 = 0

	now := time.Now().Unix()

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "|")

		if len(parts) == 2 && parts[0] == "DEL" {
			delete(m.ColdIndex, parts[1])
			offset += int64(len(line) + 1)
			continue
		}

		if len(parts) < 3 {
			offset += int64(len(line) + 1)
			continue
		}

		key := parts[0]

		expiry, err := strconv.ParseInt(parts[2], 10, 64)

		if err != nil {
			offset += int64(len(line) + 1)
			continue
		}

		if expiry != 0 && expiry < now {
			offset += int64(len(line) + 1)
			continue
		}

		m.ColdIndex[key] = ColdEntry{
			Offset: offset,
			Expiry: expiry,
		}
		offset += int64(len(line) + 1)

	}
	return scanner.Err()
}
