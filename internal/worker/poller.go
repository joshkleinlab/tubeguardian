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

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Stopping poller...")
			return
		default:
			p.processComments(ctx)
			time.Sleep(5 * time.Minute)
		}
	}
}

// processComments fetches and filters comments
func (p *Poller) processComments(ctx context.Context) {
	var comments []youtube.Comment
	var err error

	state := p.client.LoadState()

	if state.Mode == "init" {
		// First run → backfill
		log.Println("📥 First run detected → Performing full backfill scan of channel comments...")
		comments, err = p.client.FetchAllComments(ctx)

		if err == nil && len(comments) > 0 {
			newLastID := comments[0].ID
			p.client.SaveState(youtube.State{
				Mode:   "backfillDone",
				LastID: newLastID,
			})
			log.Printf("📌 Backfill completed. Last ID saved: %s\n", newLastID)
		}
	} else {
		// Incremental mode
		log.Println("🔄 Incremental mode → Fetching latest comments...")
		comments, err = p.client.FetchLatestComments(ctx, 50)
	}

	// if err != nil {
	// 	log.Printf("❌ Failed to fetch comments: %v\n", err)
	// 	return
	// }

	if len(comments) == 0 {
		log.Println("✅ No new comments found.")
		return
	}

	log.Printf("📌 Processing %d comments...\n", len(comments))
	var toHide []string
	for _, c := range comments {
		matches := p.f.Match(c.Text)
		if len(matches) > 0 {
			log.Printf("🚫 Blocked comment [%s]: \"%s\" | matches: %v\n", c.ID, c.Text, matches)
			toHide = append(toHide, c.ID)
		}
	}

	if len(toHide) > 0 {
		if err := p.client.HideComments(ctx, toHide); err != nil {
			log.Printf("❌ Failed to hide comments: %v", err)
		} else {
			log.Printf("✅ Hidden %d comments.", len(toHide))
		}
	} else {
		log.Println("👍 No banned keywords detected.")
	}
}
