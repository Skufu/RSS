package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
)

// HandlerUnfollow handles the unfollow command which removes a follow relationship between the current user and a feed
func HandlerUnfollow(s *app.State, cmd app.Command, user database.User) error {
	// Check if we have the right number of arguments
	if len(cmd.Args) < 1 {
		return errors.New("unfollow command requires a url argument")
	}

	url := cmd.Args[0]

	// Verify the feed exists
	ctx := context.Background()
	feed, err := s.Db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to find feed with URL %s: %w", url, err)
	}

	// Delete the feed follow record
	err = s.Db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
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
