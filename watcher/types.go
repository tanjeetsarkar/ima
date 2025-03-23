package watcher

import (
	"context"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileEvent represents a structured filesystem event
type FileEvent struct {
	Op         fsnotify.Op // Operation type
	Name       string      // Full path of the file/dir
	Size       int64       // File size (if applicable)
	Timestamp  time.Time   // Event timestamp
	IsDir      bool        // Whether it's a directory
	Additional interface{} // For custom metadata
}

// Config holds watcher configuration
type Config struct {
	Path         string
	EventBuffer  int
	ExcludeDirs  []string
	IncludeTypes []string
}

// Watcher interface defines the contract
type Watcher interface {
	Start(ctx context.Context) <-chan FileEvent
	Stop() error
}

// FSWatcher implements Watcher using fsnotify
type FSWatcher struct {
	config  Config
	watcher *fsnotify.Watcher
	ctx     context.Context
	cancel  context.CancelFunc
}
