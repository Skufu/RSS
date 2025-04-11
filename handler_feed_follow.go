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
func handlerFollow(s *state, cmd command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.args) < 1 {
		return errors.New("follow command requires a url argument")
	}

	url := cmd.args[0]

	// Get the feed by URL
	ctx := context.Background()
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

// handlerUnfollow handles the unfollow command which removes a follow relationship between the current user and a feed
func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.args) < 1 {
		return errors.New("unfollow command requires a url argument")
	}

	url := cmd.args[0]

	// Verify the feed exists
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to find feed with URL %s: %w", url, err)
	}

	// Delete the feed follow record
	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	// Print success message
	fmt.Printf("You (%s) have unfollowed the feed: %s\n", user.Name, feed.Name)

	return nil
}

// handlerFollowing handles the following command which lists all feeds that the current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	// Get all feed follows for the user
	ctx := context.Background()
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
