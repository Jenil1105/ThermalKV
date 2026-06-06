package persistence

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"thermalkv/internal/model"
)

func WriteLog(operation string, key string, value ...string) {
	file, err := os.OpenFile(
		"data/wal.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Println("err opening WAL", err)
		return
	}

	defer file.Close()

	log := fmt.Sprintf("%s %s %s\n", operation, key, strings.Join(value, " "))
	file.WriteString(log)

}

func LoadLogsFromDir(dir string) []string {
	files, err := os.ReadDir(dir)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var logs []string

	for _, file := range files {
		name := file.Name()

		if strings.HasPrefix(name, "wal_") && strings.HasSuffix(name, ".log") {
			path := filepath.Join(dir, name)

			f, err := os.Open(path)

			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				logs = append(logs, scanner.Text())
			}
			f.Close()

		}
	}
	file, err := os.Open(filepath.Join(dir, "wal.log"))

	if err != nil {
		return logs
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	return logs
}

func LoadLogs() []string {
	return LoadLogsFromDir("data")
}

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

func LoadSnapshot() map[string]model.SnapshotItem {
	file, err := os.Open("data/snapshot.dat")

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

func ClearWAL() {
	file, err := os.Create("data/wal.log")

	if err != nil {
		fmt.Println("Error clearing WAL:", err)
		return
	}

	file.Close()

}
