package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Skufu/RSS/internal/database"
	"github.com/google/uuid"
)

// handlerAddFeed handles the addfeed command which adds a feed for the current user
func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.args) < 2 {
		return errors.New("addfeed command requires name and url arguments")
	}

	name := cmd.args[0]
	url := cmd.args[1]

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

	feed, err := s.db.CreateFeed(ctx, feedParams)
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

	feedFollow, err := s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("feed created but failed to follow it: %w", err)
	}

	fmt.Printf("You (%s) are now following this feed.\n", feedFollow.UserName)

	return nil
}

// handlerFeeds handles the feeds command which lists all feeds in the database
func handlerFeeds(s *state, cmd command) error {
	// Get all feeds from the database
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}

	// If no feeds, print a message
	if len(feeds) == 0 {
		fmt.Println("No feeds found in the database.")
		return nil
	}

	// Print header
	fmt.Println("Feeds:")
	fmt.Println("------")

	// Print each feed with its user
	for i, feed := range feeds {
		fmt.Printf("%d. %s\n", i+1, feed.Name)
		fmt.Printf("   URL: %s\n", feed.Url)
		fmt.Printf("   Added by: %s\n", feed.UserName)
		fmt.Println()
	}

	return nil
}
