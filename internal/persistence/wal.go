package persistence

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type WAL struct {
	File         *os.File
	Writer       *bufio.Writer
	Sync         bool
	SyncInterval time.Duration
	stopChan     chan struct{}
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

	writer := bufio.NewWriterSize(file, 64*1024)

	return &WAL{
		File:         file,
		Sync:         sync,
		Writer:       writer,
		SyncInterval: 100 * time.Millisecond,
		stopChan:     make(chan struct{}),
	}
}

func (w *WAL) Write(operation string, key string, value ...string) error {

	var builder strings.Builder

	builder.WriteString(operation)
	builder.WriteByte(' ')
	builder.WriteString(key)

	for _, v := range value {
		builder.WriteByte(' ')
		builder.WriteString(v)
	}

	builder.WriteByte('\n')

	_, err := w.Writer.WriteString(builder.String())

	return err
}

func (w *WAL) Close() {
	if w != nil {
		if w.Sync {
			close(w.stopChan)
		}
		if w.Writer != nil {
			w.Writer.Flush()
		}
		if w.File != nil {
			w.File.Close()
		}
	}
}

func (w *WAL) StartSyncLoop() {

	if !w.Sync {
		return
	}

	ticker := time.NewTicker(w.SyncInterval)

	go func() {

		for {
			select {

			case <-ticker.C:

				w.Writer.Flush()
				w.File.Sync()

			case <-w.stopChan:

				ticker.Stop()

				w.Writer.Flush()
				w.File.Sync()

				return
			}
		}
	}()
}
