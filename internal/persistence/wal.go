package persistence

import (
	"fmt"
	"os"
	"strings"
)

type WAL struct {
	File *os.File
	Sync bool
}

func NewWAL(sync bool) *WAL {
	file, err := os.OpenFile(
		"data/wal.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Println("WAL open error: ", err)
		return nil
	}
	return &WAL{
		File: file,
		Sync: sync,
	}
}

func (w *WAL) Write(operation string, key string, value ...string) error {
	log := fmt.Sprintf("%s %s %s\n", operation, key, strings.Join(value, " "))
	_, err := w.File.WriteString(log)

	if err != nil {
		return err
	}

	if w.Sync {
		err = w.File.Sync()
		if err != nil {
			return err
		}
	}
	return nil

}

func (w *WAL) Close() {
	if w != nil && w.File != nil {
		w.File.Close()
	}
}
