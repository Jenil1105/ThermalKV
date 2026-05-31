package persistence

import (
	"fmt"
	"os"
	"strings"
)

type WAL struct {
	File *os.File
}

func NewWAL() *WAL {
	file, err := os.OpenFile(
		"data/wal.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Println("WAL open error: ", err)
		return nil
	}
	return &WAL{File: file}
}

func (w *WAL) Write(operation string, key string, value ...string) {
	log := fmt.Sprintf("%s %s %s\n", operation, key, strings.Join(value, " "))
	w.File.WriteString(log)
	w.File.Sync()
}

func (w *WAL) Close() {
	if w != nil && w.File != nil {
		w.File.Close()
	}
}
