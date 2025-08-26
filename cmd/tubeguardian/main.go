package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joshkleinlab/tubeguardian/internal/config"
	"github.com/joshkleinlab/tubeguardian/internal/filter"
	"github.com/joshkleinlab/tubeguardian/internal/worker"
	"github.com/joshkleinlab/tubeguardian/internal/youtube"
)

func main() {
	// Load config.yaml
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Load banned words
	matcher, err := filter.LoadKeywords(cfg.BannedWordsFile)
	if err != nil {
		log.Fatalf("❌ Failed to load banned words: %v", err)
	}

	// Initialize YouTube client
	ctx := context.Background()
	client, err := youtube.NewClient(ctx, cfg.YourChannelID, "configs/credentials.json")
	if err != nil {
		log.Fatalf("❌ Failed to create YouTube client: %v", err)
	}

	// Setup poller
	p := worker.NewPoller(client, matcher, cfg)

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Capture SIGINT / SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("🛑 Received shutdown signal. Cleaning up...")
		cancel()
	}()

	// Start the poller
	p.Run(ctx)
}
