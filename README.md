# Gator - A Command-Line RSS Feed Aggregator

Gator is a command-line RSS feed aggregator that lets you:
- Follow your favorite RSS feeds
- Automatically collect and store posts
- Browse posts from the terminal

## Prerequisites

To run Gator, you'll need:

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) (version 12 or higher)

## Installation

### Installing from Source

```bash
# Clone the repository
git clone https://github.com/Skufu/RSS.git
cd RSS

# Install the binary
go install
```

After installation, the `gator` command should be available in your PATH.

### Setting Up the Database

1. Create a PostgreSQL database:

```bash
createdb rss
```

2. Run the migrations:

```bash
# Install goose migration tool if you don't have it
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/schema postgres "postgres://localhost:5432/rss?sslmode=disable" up
```

## Configuration

Gator uses a configuration file located at `~/.gatorconfig.json`:

```json
{
  "current_user_name": "",
  "database_url": "postgres://localhost:5432/rss?sslmode=disable"
}
```

You can create this file manually or let Gator create it with default values.

## Commands

Here are some of the commands you can use:

### User Management

```bash
# Register a new user
gator register username

# Login as a user
gator login username
```

### Feed Management

```bash
# Add a new feed
gator addfeed "HackerNews" "https://news.ycombinator.com/rss"

# List all feeds
gator feeds

# Follow an existing feed
gator follow "https://news.ycombinator.com/rss"

# List feeds you're following
gator following

# Unfollow a feed
gator unfollow "https://news.ycombinator.com/rss"
```

### Content Aggregation

```bash
# Start the aggregator (runs continuously, collecting posts every 5 minutes)
gator agg 5m

# Browse recent posts from feeds you follow (default 20 posts)
gator browse

# Browse with custom limit (e.g., only 5 posts)
gator browse 5
```

## How It Works

1. **Register**: Create a user account
2. **Add Feeds**: Add RSS feeds you want to follow
3. **Follow Feeds**: Choose which feeds to follow
4. **Aggregate**: Run the aggregator to collect posts
5. **Browse**: View posts from feeds you follow

The aggregator runs in the background and continuously collects new posts from your followed feeds at the interval you specify.

## Development

For development, you can use:

```bash
go run . <command>
```

But for production use, install the binary and use:

```bash
gator <command>
```

## License

MIT 