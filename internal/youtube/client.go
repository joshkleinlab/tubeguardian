package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	youtube "google.golang.org/api/youtube/v3"
)

// Client wraps YouTube API service + state
type Client struct {
	channelID string
	service   *youtube.Service
	stateFile string // stores last processed comment ID
}

// Comment represents a YouTube comment
type Comment struct {
	ID   string
	Text string
}

// NewClient initializes YouTube client with OAuth2
func NewClient(ctx context.Context, channelID, credentialsFile string) (*Client, error) {
	service, err := NewYouTubeService(ctx, credentialsFile)
	if err != nil {
		return nil, err
	}

	stateFile := filepath.Join(".", "configs", "state.json")

	return &Client{
		channelID: channelID,
		service:   service,
		stateFile: stateFile,
	}, nil
}

// FetchAllComments performs a global scan (first run) with paging
func (c *Client) FetchAllComments(ctx context.Context) ([]Comment, error) {
	var allComments []Comment
	nextPageToken := ""

	for {
		call := c.service.CommentThreads.List([]string{"snippet"}).
			AllThreadsRelatedToChannelId(c.channelID).
			Order("time").
			MaxResults(100).
			PageToken(nextPageToken)

		resp, err := call.Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("API error (FetchAllComments): %w", err)
		}

		for _, item := range resp.Items {
			snippet := item.Snippet.TopLevelComment.Snippet
			allComments = append(allComments, Comment{
				ID:   item.Snippet.TopLevelComment.Id,
				Text: snippet.TextDisplay,
			})
		}

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}
	return allComments, nil
}

// FetchLatestComments fetches recent comments
// FetchLatestComments gets only new comments since last run
func (c *Client) FetchLatestComments(ctx context.Context, maxResults int64) ([]Comment, error) {
	// Load previous state
	state := c.LoadState()
	lastID := state.LastID

	call := c.service.CommentThreads.List([]string{"snippet"}).
		ChannelId(c.channelID).
		Order("time").
		MaxResults(maxResults)

	resp, err := call.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("API error (FetchLatestComments): %w", err)
	}

	var comments []Comment
	for _, item := range resp.Items {
		cmt := item.Snippet.TopLevelComment.Snippet
		commentID := item.Snippet.TopLevelComment.Id

		// Stop when reaching the last seen comment
		if commentID == lastID {
			break
		}

		comments = append(comments, Comment{
			ID:   commentID,
			Text: cmt.TextDisplay,
		})
	}

	// Save newest comment ID for next run
	if len(resp.Items) > 0 {
		newLastID := resp.Items[0].Snippet.TopLevelComment.Id
		c.SaveState(State{
			Mode:   "backfillDone",
			LastID: newLastID,
		})
	}

	return comments, nil
}

// HideComment hides a single comment
func (c *Client) HideComment(ctx context.Context, commentID string) error {
	return c.HideComments(ctx, []string{commentID})
}

// HideComments hides multiple comments at once
func (c *Client) HideComments(ctx context.Context, commentIDs []string) error {
	if len(commentIDs) == 0 {
		return nil
	}

	call := c.service.Comments.SetModerationStatus(commentIDs, "heldForReview")
	err := call.Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to hide comments %v: %w", commentIDs, err)
	}
	return nil
}

type State struct {
	Mode   string `json:"mode"`    // "init" (before backfill) or "backfillDone"
	LastID string `json:"last_id"` // most recent processed comment ID
}

// loadState returns the saved state, or default if none
func (c *Client) LoadState() State {
	f, err := os.Open(c.stateFile)
	if err != nil {
		return State{Mode: "init", LastID: ""}
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return State{Mode: "init", LastID: ""}
	}
	return s
}

// saveState persists state to file
func (c *Client) SaveState(s State) {
	f, err := os.Create(c.stateFile)
	if err != nil {
		fmt.Printf("⚠️ Unable to save state: %v\n", err)
		return
	}
	defer f.Close()

	json.NewEncoder(f).Encode(s)
}
