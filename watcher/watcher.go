package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// New creates a new file watcher instance
func New(cfg Config) (*FSWatcher, error) {
	cleanPath := filepath.Clean(cfg.Path)
	if cfg.EventBuffer <= 0 {
		cfg.EventBuffer = 100
	}

	stat, err := os.Stat(cleanPath)
	if err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("invalid directory path: %w", err)
	}

	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &FSWatcher{
		config:  cfg,
		watcher: fswatcher,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Start begins watching for file system events
func (w *FSWatcher) Start(ctx context.Context) <-chan FileEvent {
	eventChan := make(chan FileEvent, w.config.EventBuffer)

	// Add initial directories
	err := filepath.WalkDir(w.config.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if shouldIgnore(path, w.config.ExcludeDirs) {
			if d.IsDir() {
				return filepath.SkipDir // Skip the entire directory
			}
			return nil // Skip the file
		}
		if d.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Initial directory walk error: %v", err)
	}

	go w.processEvents(ctx, eventChan)
	return eventChan
}

// Stop shuts down the watcher
func (w *FSWatcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *FSWatcher) processEvents(ctx context.Context, eventChan chan<- FileEvent) {
	defer close(eventChan)

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.ctx.Done():
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event, eventChan)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (w *FSWatcher) handleEvent(event fsnotify.Event, eventChan chan<- FileEvent) {

	if shouldIgnore(event.Name, w.config.ExcludeDirs) {
		return
	}

	// Handle directory creation
	if event.Op.Has(fsnotify.Create) {
		info, err := os.Stat(event.Name)
		if err == nil && info.IsDir() {
			_ = w.watcher.Add(event.Name)
		}
	}

	// Skip unwanted events
	if event.Op.Has(fsnotify.Chmod) {
		return
	}

	// Collect file info
	var size int64
	var isDir bool
	if info, err := os.Stat(event.Name); err == nil {
		size = info.Size()
		isDir = info.IsDir()
	}

	// Create structured event
	fileEvent := FileEvent{
		Op:        event.Op,
		Name:      event.Name,
		Size:      size,
		Timestamp: time.Now().UTC(),
		IsDir:     isDir,
	}

	select {
	case <-w.ctx.Done():
		return
	case eventChan <- fileEvent:
	}
}

func shouldIgnore(path string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}
	}
	return false
}
