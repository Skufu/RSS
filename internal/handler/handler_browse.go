package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
)

// formatTime returns a human-readable string for a time
func formatTime(t time.Time) string {
	return t.Format("Jan 02, 2006")
}

// HandlerBrowse handles the browse command which displays posts from feeds the user is following
func HandlerBrowse(s *app.State, cmd app.Command, user database.User) error {
	// Default limit to 20 if not provided
	limit := int32(20)
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit parameter: %w", err)
		}
		limit = int32(parsedLimit)
	}

	// Get posts for the user
	ctx := context.Background()
	posts, err := s.Db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	// Check if there are any posts
	if len(posts) == 0 {
		fmt.Println("No posts found. Try following some feeds and waiting for the aggregator to collect posts.")
		return nil
	}

	// Print header
	fmt.Printf("Latest %d posts from your feeds:\n", len(posts))
	fmt.Println("==================================")
	fmt.Println()

	// Print each post
	for i, post := range posts {
		// Add a divider between posts except for the first one
		if i > 0 {
			fmt.Println("----------------------------------")
		}

		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Feed: %s\n", post.FeedName)

		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", formatTime(post.PublishedAt.Time))
		}

		fmt.Printf("URL: %s\n", post.Url)

		if post.Description.Valid && post.Description.String != "" {
			// Print first 200 chars of description or less
			desc := post.Description.String
			if len(desc) > 200 {
				desc = desc[:197] + "..."
			}
			fmt.Printf("Description: %s\n", desc)
		}

		fmt.Println()
	}

	return nil
}
