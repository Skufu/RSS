package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Skufu/RSS/internal/database"
	"github.com/google/uuid"
)

// handlerFollow handles the follow command which creates a follow relationship between the current user and a feed
func handlerFollow(s *state, cmd command) error {
	// Check if we have the right number of arguments
	if len(cmd.args) < 1 {
		return errors.New("follow command requires a url argument")
	}

	url := cmd.args[0]

	// Check if a user is set in the config
	if s.cfg.CurrentUserName == "" {
		return errors.New("no user set, please login first")
	}

	// Get the current user from the database
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Get the feed by URL
	feed, err := s.db.GetFeedByURL(ctx, url)
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

	feedFollow, err := s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	// Print success message
	fmt.Printf("You (%s) are now following the feed: %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

// handlerFollowing handles the following command which lists all feeds that the current user is following
func handlerFollowing(s *state, cmd command) error {
	// Check if a user is set in the config
	if s.cfg.CurrentUserName == "" {
		return errors.New("no user set, please login first")
	}

	// Get the current user from the database
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Get all feed follows for the user
	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	// If no feeds followed, print a message
	if len(feedFollows) == 0 {
		fmt.Printf("You (%s) are not following any feeds.\n", s.cfg.CurrentUserName)
		return nil
	}

	// Print header
	fmt.Printf("Feeds followed by %s:\n", s.cfg.CurrentUserName)
	fmt.Println("--------------------------------------")

	// Print each followed feed
	for i, ff := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, ff.FeedName)
	}

	return nil
}
