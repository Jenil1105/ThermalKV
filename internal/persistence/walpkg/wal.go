package walpkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type WAL struct {
	Path         string
	File         *os.File
	Writer       *bufio.Writer
	Sync         bool
	SyncInterval time.Duration
	stopChan     chan struct{}

	Mutex sync.Mutex
}

func NewWAL(path string, sync bool) *WAL {

	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Println("WAL open error: ", err)
		return nil
	}

	writer := bufio.NewWriterSize(file, 64*1024)

	return &WAL{
		Path:         path,
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

	dir := filepath.Dir(w.Path)
	base := filepath.Base(w.Path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	rotatedFile := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext))

	err = os.Rename(w.Path, rotatedFile)
	if err != nil {
		return "", err
	}

	newfile, err := os.OpenFile(
		w.Path,
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
