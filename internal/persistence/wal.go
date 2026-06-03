package persistence

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type WAL struct {
	File         *os.File
	Writer       *bufio.Writer
	Sync         bool
	SyncInterval time.Duration
	stopChan     chan struct{}

	Mutex sync.Mutex
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

	fmt.Println("WAL write called")

	var builder strings.Builder

	builder.WriteString(operation)
	builder.WriteByte(' ')
	builder.WriteString(key)

	for _, v := range value {
		builder.WriteByte(' ')
		builder.WriteString(v)
	}

	builder.WriteByte('\n')

	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	_, err := w.Writer.WriteString(builder.String())
	if err != nil {
		return err
	}

	if !w.Sync {
		return w.Writer.Flush()
	}

	return nil
}

func (w *WAL) Close() {
	fmt.Println("WAL CLOSED")
	if w != nil {
		if w.Sync {
			close(w.stopChan)
		}

		w.Mutex.Lock()
		defer w.Mutex.Unlock()

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

				w.Mutex.Lock()

				w.Writer.Flush()
				w.File.Sync()

				w.Mutex.Unlock()

			case <-w.stopChan:

				ticker.Stop()

				w.Mutex.Lock()

				w.Writer.Flush()
				w.File.Sync()

				w.Mutex.Unlock()

				return
			}
		}
	}()
}

func (w *WAL) Rotate() (string, error) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	err := w.Writer.Flush()

	if err != nil {
		return "", err
	}

	err = w.File.Close()
	if err != nil {
		return "", err
	}

	rotatedFile := fmt.Sprintf("data/wal_%d.log", time.Now().UnixNano())

	err = os.Rename("data/wal.log", rotatedFile)
	if err != nil {
		return "", err
	}

	newfile, err := os.OpenFile(
		"data/wal.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return "", err
	}

	w.File = newfile
	w.Writer = bufio.NewWriterSize(newfile, 64*1024)

	return rotatedFile, nil
}
