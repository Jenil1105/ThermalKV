package walpkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"strings"
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

func LoadLogs(dir string) []string {
	return LoadLogsFromDir(dir)
}

func ClearWAL() {
	file, err := os.Create("data/wal.log")

	if err != nil {
		fmt.Println("Error clearing WAL:", err)
		return
	}

	file.Close()

}
