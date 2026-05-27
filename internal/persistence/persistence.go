package persistence

import (
	"bufio"
	"fmt"
	"os"
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

func LoadLogs() []string {
	file, err := os.Open("data/wal.log")

	if err != nil {
		fmt.Println("err opening file", err)
		return nil
	}
	defer file.Close()
	var logs []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	return logs
}

func SaveSnapshot(data map[string]model.SnapshotItem) {

	file, err := os.Create("data/snapshot.dat")

	if err != nil {
		fmt.Println("Snapshot err", err)
		return
	}

	defer file.Close()

	for key, item := range data {
		line := fmt.Sprintf("%s|%s|%d\n", key, item.Value, item.Expiry)
		file.WriteString(line)
	}

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
