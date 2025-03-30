package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/image-archive/indexer"
	"github.com/image-archive/watcher"
	"github.com/spf13/cobra"
)

func main() {
	var watchFlag bool

	rootCmd := &cobra.Command{
		Use:   "image-archive [directory]",
		Short: "Image Archive CLI",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Error: Directory path is required")
			}

			dir := args[0]

			indexer.SplitCreate(dir)

			if watchFlag {
				log.Println("Starting watcher...")
				watcherCmd(dir)
			}
		},
	}
	indexer.AddFlags(rootCmd)
	rootCmd.Flags().BoolVar(&watchFlag, "watch", false, "Start watching the directory for changes")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func watcherCmd(dir string) {
	cfg := watcher.Config{
		Path:        dir,
		EventBuffer: 100,
		ExcludeDirs: []string{"index.html", ".thumbs"},
	}

	fileWatcher, err := watcher.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer fileWatcher.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	eventChan := fileWatcher.Start(ctx)

	go watcher.EventConsumer(ctx, eventChan)

	log.Printf("Watching Directory: %s (PID: %d)", dir, os.Getpid())

	select {
	case <-sigChan:
		log.Println("Shutdown signal received")
	case <-ctx.Done():
	}
}
