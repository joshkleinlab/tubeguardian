package worker

import (
	"context"
	"log"
	"time"

	"github.com/joshkleinlab/tubeguardian/internal/config"
	"github.com/joshkleinlab/tubeguardian/internal/filter"
	"github.com/joshkleinlab/tubeguardian/internal/youtube"
)

// Poller runs comment fetching & filtering
type Poller struct {
	client *youtube.Client
	f      *filter.Matcher
	cfg    *config.Config
}

// NewPoller creates a new poller
func NewPoller(client *youtube.Client, f *filter.Matcher, cfg *config.Config) *Poller {
	return &Poller{client: client, f: f, cfg: cfg}
}

// Run starts periodic comment fetching and filtering
func (p *Poller) Run(ctx context.Context) {
	log.Println("🚀 TubeGuardian started. Press Ctrl+C to stop.")
	state := p.client.LoadState()

	// If first run → do a full backfill once
	if state.Mode == "init" {
		log.Println("📥 First run → Performing full backfill...")
		if err := p.client.FetchAllComments(ctx); err != nil {
			log.Printf("❌ Backfill failed: %v", err)
		}
		_ = p.client.SaveState(youtube.State{
			Mode:   "backfillDone",
			LastID: "", // will be updated as comments stream
		})
	}

	// Start comment consumer
	go p.consumeComments(ctx)

	// Poll new comments every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Poller stopped.")
			return
		case <-ticker.C:
			log.Println("🔄 Fetching latest comments...")
			if err := p.client.FetchLatestComments(ctx, 50); err != nil {
				log.Printf("❌ Failed to fetch latest comments: %v", err)
			}
		}
	}
}

// consumeComments
func (p *Poller) consumeComments(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-p.client.Out:
			matches := p.f.Match(c.Text)
			if len(matches) > 0 {
				log.Printf("🚫 Blocked [%s]: \"%s\" | matches: %v", c.ID, c.Text, matches)
				if err := p.client.HideComments(ctx, []string{c.ID}, p.cfg.ModeRation); err != nil {
					log.Printf("❌ Failed to hide comment %s: %v", c.ID, err)
				} else {
					log.Printf("✅ Hidden comment %s", c.ID)
				}
			}

			// Save latest ID
			_ = p.client.SaveState(youtube.State{
				Mode:   "backfillDone",
				LastID: c.ID,
			})
		}
	}
}
