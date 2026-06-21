package fs

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	whlogger "github.com/bohdanch-w/wheel/logger"
)

func NewHotSwapFileReader(logger whlogger.Logger, filePath string) (*HotSwapFileReader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create file watcher: %w", err)
	}

	file := &HotSwapFileReader{
		logger:      logger,
		filePath:    filePath,
		done:        make(chan struct{}),
		fileWatcher: watcher,
	}

	if err := file.start(); err != nil {
		return nil, fmt.Errorf("start file reader: %w", err)
	}

	return file, nil
}

type HotSwapFileReader struct {
	logger      whlogger.Logger
	fileWatcher *fsnotify.Watcher

	filePath string
	content  []byte
	changed  bool

	mu   sync.RWMutex
	done chan struct{}
}

func (r *HotSwapFileReader) start() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("get file data: %w", err)
	}

	r.content = data
	r.changed = true

	if err := r.fileWatcher.Add(r.filePath); err != nil {
		return fmt.Errorf("add file to watcher: %w", err)
	}

	go r.watchFile()

	return nil
}

func (r *HotSwapFileReader) watchFile() {
	// os sends multiple events for the same file change
	// trigger content update only after last event received
	// new events in the same time window will reset the timer
	const eventCooldown = time.Millisecond * 100

	ticker := time.NewTicker(eventCooldown)
	ticker.Stop()

	for {
		select {
		case <-r.done:
			return
		case event := <-r.fileWatcher.Events:
			if event.Has(fsnotify.Write) {
				r.logger.Infof("file updated: %q", r.filePath)
				ticker.Reset(eventCooldown)
			}
		case <-ticker.C:
			data, err := os.ReadFile(r.filePath)
			if err != nil {
				r.logger.Errorf("update file %q: %w", r.filePath, err)

				continue
			}

			r.mu.Lock()
			r.content = data
			r.changed = true
			r.mu.Unlock()

			ticker.Stop()

			r.logger.Infof("content updated: %q", r.filePath)
		}
	}
}

func (r *HotSwapFileReader) Content() (_ []byte, changed bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	changed = r.changed
	r.changed = false

	return r.content, changed
}

func (r *HotSwapFileReader) Close() {
	close(r.done)
	r.fileWatcher.Close()
}
