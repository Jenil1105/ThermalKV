package snapshot

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"thermalkv/internal/model"
)

func SaveSnapshot(data map[string]model.SnapshotItem) error {

	file, err := os.Create("data/snapshot.dat")

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for key, item := range data {
		line := fmt.Sprintf("%s|%s|%d\n", key, item.Value, item.Expiry)
		_, err := writer.WriteString(line)

		if err != nil {
			return err
		}
	}
	return nil

}

func LoadSnapshot(path string) map[string]model.SnapshotItem {
	file, err := os.Open(path)

	if err != nil {
		return nil
	}

	defer file.Close()

	snapshot := make(map[string]model.SnapshotItem)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "|")

		if len(parts) < 3 {
			continue
		}

		key := parts[0]
		value := parts[1]

		expiry, err := strconv.ParseInt(parts[2], 10, 64)

		if err != nil {
			continue
		}

		snapshot[key] = model.SnapshotItem{
			Value:  value,
			Expiry: expiry,
		}
	}
	return snapshot
}
