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
	}, true

}
