package recover

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"thermalkv/internal/coldstore"
	"time"
)

func RecoverColdIndex(m *coldstore.Manager) error {

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

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

		m.ColdIndex[key] = coldstore.ColdEntry{
			Offset: offset,
			Expiry: expiry,
		}
		offset += int64(len(line) + 1)

	}
	return scanner.Err()
}
