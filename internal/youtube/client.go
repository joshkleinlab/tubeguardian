package youtube

import (
	"context"
	"fmt"

	youtube "google.golang.org/api/youtube/v3"
)

// Client wraps YouTube API service + state
type Client struct {
	channelID string
	service   *youtube.Service
	Out       chan Comment // 🔑 Channel for streaming comments
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

	return &Client{
		channelID: channelID,
		service:   service,
		Out:       make(chan Comment, 150), // buffered
	}, nil
}

// FetchAllComments performs a global scan (first run) with paging
func (c *Client) FetchAllComments(ctx context.Context) error {
	call := c.service.CommentThreads.List([]string{"snippet"}).
		AllThreadsRelatedToChannelId(c.channelID).
		Order("time").
		MaxResults(100)

	for {
		resp, err := call.Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("API error (FetchAllComments): %w", err)
		}

		for _, item := range resp.Items {
			cmt := item.Snippet.TopLevelComment.Snippet
			c.Out <- Comment{
				ID:   item.Snippet.TopLevelComment.Id,
				Text: cmt.TextDisplay,
			}
		}

		if resp.NextPageToken == "" {
			break
		}
		call = call.PageToken(resp.NextPageToken)
	}

	return nil
}

// FetchLatestComments gets only new comments since last run
func (c *Client) FetchLatestComments(ctx context.Context, maxResults int64) error {
	state := c.LoadState()

	call := c.service.CommentThreads.List([]string{"snippet"}).
		ChannelId(c.channelID).
		Order("time").
		MaxResults(maxResults)

	resp, err := call.Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("API error (FetchLatestComments): %w", err)
	}

	for _, item := range resp.Items {
		cmt := item.Snippet.TopLevelComment.Snippet
		comment := Comment{
			ID:   item.Snippet.TopLevelComment.Id,
			Text: cmt.TextDisplay,
		}

		if comment.ID == state.LastID {
			break
		}

		c.Out <- comment
	}

	return nil
}

// HideComment hides a single comment
func (c *Client) HideComment(ctx context.Context, commentID string) error {
	return c.HideComments(ctx, []string{commentID}, "heldForReview")
}

// HideComments hides multiple comments at once
func (c *Client) HideComments(ctx context.Context, commentIDs []string, moderationStatus string) error {
	if len(commentIDs) == 0 {
		return nil
	}

	call := c.service.Comments.SetModerationStatus(commentIDs, moderationStatus)
	err := call.Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to hide comments %v: %w", commentIDs, err)
	}
	return nil
}
