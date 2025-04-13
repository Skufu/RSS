package handler

import (
	"context"
	"fmt"

	"github.com/Skufu/RSS/internal/app"
)

// HandlerFeeds handles the feeds command which lists all feeds in the database
func HandlerFeeds(s *app.State, cmd app.Command) error {
	// Get all feeds from the database
	ctx := context.Background()
	feeds, err := s.Db.GetFeeds(ctx)
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
