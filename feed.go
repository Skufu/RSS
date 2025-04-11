package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

// RSSFeed represents the structure of an RSS feed
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents an item in an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed retrieves and parses an RSS feed from the given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set User-Agent header to identify our client
	req.Header.Set("User-Agent", "gator")

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal XML into RSSFeed struct
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %w", err)
	}

	// Decode HTML entities in channel fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Decode HTML entities in item fields
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

// handlerAgg handles the agg command which fetches and displays a feed
func handlerAgg(s *state, cmd command) error {
	// For now, we'll just fetch the wagslane.dev feed
	feedURL := "https://www.wagslane.dev/index.xml"

	fmt.Printf("Fetching feed from %s\n", feedURL)

	ctx := context.Background()
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	// Print the feed information
	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	fmt.Printf("Items: %d\n\n", len(feed.Channel.Item))

	// Print each item
	for i, item := range feed.Channel.Item {
		fmt.Printf("Item %d:\n", i+1)
		fmt.Printf("  Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Publication Date: %s\n", item.PubDate)
		fmt.Printf("  Description: %s\n\n", item.Description)
	}

	return nil
}
