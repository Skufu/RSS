package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
	"github.com/Skufu/RSS/internal/feed"
	"github.com/google/uuid"
)

// parsePublishedAt tries to parse the published date in various common formats
func parsePublishedAt(dateStr string) (time.Time, error) {
	// Common RSS date formats
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// scrapeFeeds is a helper function that fetches a single feed
func scrapeFeeds(s *app.State) error {
	ctx := context.Background()

	// Get the next feed to fetch
	feedItem, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}

	fmt.Printf("Fetching feed: %s (%s)\n", feedItem.Name, feedItem.Url)

	// Mark feed as fetched first to avoid issues with failed fetches
	err = s.Db.MarkFeedFetched(ctx, feedItem.ID)
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	// Fetch the feed
	feedData, err := feed.FetchFeed(ctx, feedItem.Url)
	if err != nil {
		// Just log the error and return nil so we continue with other feeds
		fmt.Printf("Error: failed to fetch feed %s: %v\n", feedItem.Name, err)
		return nil
	}

	// Print the feed information
	fmt.Printf("Feed: %s\n", feedItem.Name)
	fmt.Printf("Items: %d\n", len(feedData.Channel.Item))

	// Save each post to the database
	for _, item := range feedData.Channel.Item {
		// Try to parse the published date
		var publishedAt sql.NullTime
		if item.PubDate != "" {
			if parsedTime, err := parsePublishedAt(item.PubDate); err == nil {
				publishedAt = sql.NullTime{
					Time:  parsedTime,
					Valid: true,
				}
			} else {
				// If parsing fails, just log and continue with a null time
				fmt.Printf("Warning: Could not parse date '%s' for post '%s': %v\n",
					item.PubDate, item.Title, err)
			}
		}

		// Prepare post parameters
		now := time.Now()
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       strings.TrimSpace(item.Title),
			Url:         strings.TrimSpace(item.Link),
			Description: sql.NullString{String: strings.TrimSpace(item.Description), Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      feedItem.ID,
		}

		// Create the post, ignoring duplicate URL errors
		_, err := s.Db.CreatePost(ctx, postParams)
		if err != nil {
			// Check if it's a unique constraint violation on the URL
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
				strings.Contains(err.Error(), "posts_url_key") {
				// Silently ignore duplicate URLs
				continue
			}

			// For other errors, log but continue processing
			fmt.Printf("Error saving post '%s': %v\n", item.Title, err)
		} else {
			fmt.Printf("  → Saved: %s\n", item.Title)
		}
	}

	fmt.Println()

	return nil
}

// HandlerAgg handles the agg command which fetches and displays feeds in a loop
func HandlerAgg(s *app.State, cmd app.Command) error {
	// Check for time_between_reqs argument
	if len(cmd.Args) < 1 {
		return errors.New("agg command requires a time_between_reqs argument (e.g. 10s, 1m, 1h)")
	}

	// Parse the duration
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	fmt.Printf("Starting feed aggregator. Collecting feeds every %v. Press Ctrl+C to stop.\n", timeBetweenRequests)

	// Set up a channel to handle interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a ticker for periodic feed scraping
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	// Start the scraping loop in a goroutine
	go func() {
		// Run immediately without waiting for first tick
		for ; ; <-ticker.C {
			err := scrapeFeeds(s)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nStopping feed aggregator.")

	return nil
}
