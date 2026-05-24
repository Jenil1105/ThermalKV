package persistence

import (
	"bufio"
	"fmt"
	"os"
)

func WriteLog(operation string, key string, value string) {
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

	log := fmt.Sprintf("%s %s %s\n", operation, key, value)
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
