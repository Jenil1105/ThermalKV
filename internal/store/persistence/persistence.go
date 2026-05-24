package persistence

import (
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
