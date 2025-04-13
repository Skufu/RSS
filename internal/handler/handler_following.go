package handler

import (
	"context"
	"fmt"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
)

// HandlerFollowing handles the following command which lists all feeds that the current user is following
func HandlerFollowing(s *app.State, cmd app.Command, user database.User) error {
	// Get all feed follows for the user
	ctx := context.Background()
	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	// If no feeds followed, print a message
	if len(feedFollows) == 0 {
		fmt.Printf("You (%s) are not following any feeds.\n", s.Cfg.CurrentUserName)
		return nil
	}

	// Print header
	fmt.Printf("Feeds followed by %s:\n", s.Cfg.CurrentUserName)
	fmt.Println("--------------------------------------")

	// Print each followed feed
	for i, ff := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, ff.FeedName)
	}

	return nil
}
