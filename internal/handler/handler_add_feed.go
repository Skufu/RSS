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

// HandlerAddFeed handles the addfeed command which adds a feed for the current user
func HandlerAddFeed(s *app.State, cmd app.Command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.Args) < 2 {
		return errors.New("addfeed command requires name and url arguments")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	// Create the feed
	ctx := context.Background()
	now := time.Now()
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	// Print success message and feed data
	fmt.Printf("Feed created successfully!\n")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Created: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)

	// Create a feed follow record for the user who added the feed
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("feed created but failed to follow it: %w", err)
	}

	fmt.Printf("You (%s) are now following this feed.\n", feedFollow.UserName)

	return nil
}
