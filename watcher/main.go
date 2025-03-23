package watcher

import (
	"context"
	"log"
	"path/filepath"

	"github.com/image-archive/indexer"
)

func EventConsumer(ctx context.Context, eventChan <-chan FileEvent) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer Shutting Down")
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			processEvent(event)

		}
	}
}

func processEvent(event FileEvent) {
	log.Printf("[EVENT] %-8s %q (Size: %d, Dir: %t)",
		event.Op.String(),
		event.Name,
		event.Size,
		event.IsDir,
	)
	if event.IsDir {
		indexer.SplitCreate(event.Name)
	} else {
		indexer.SplitCreate(filepath.Dir(event.Name))
	}
}
