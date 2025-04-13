package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
	"github.com/google/uuid"
)

// HandlerFollow handles the follow command which creates a follow relationship between the current user and a feed
func HandlerFollow(s *app.State, cmd app.Command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.Args) < 1 {
		return errors.New("follow command requires a url argument")
	}

	url := cmd.Args[0]

	// Get the feed by URL
	ctx := context.Background()
	feed, err := s.Db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to find feed with URL %s: %w", url, err)
	}

	// Create the feed follow
	now := time.Now()
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	// Print success message
	fmt.Printf("You (%s) are now following the feed: %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}
