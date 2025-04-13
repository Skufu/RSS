# 📰 Gator - A Command-Line RSS Feed Aggregator

<div align="center">

![Gator RSS Logo](https://img.shields.io/badge/Gator-RSS%20Aggregator-green)
[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go&logoColor=white)](https://golang.org/doc/install)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12+-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/download/)


</div>

Gator is a command-line RSS feed aggregator that helps you stay updated with content sources from your terminal.

## Features

- **Discover & Follow**: Track your favorite RSS feeds
- **Automatic Collection**: Schedule periodic content fetching
- **Content Storage**: Save posts to a database for later access
- **Terminal Reading**: Browse all content from the command line

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go** - version 1.18 or higher [→ Installation Guide](https://golang.org/doc/install)
- **PostgreSQL** - version 12 or higher [→ Installation Guide](https://www.postgresql.org/download/)

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

### Database Setup

1. **Create a PostgreSQL database**:

```bash
createdb rss
```

2. **Run the migrations**:

```bash
# Install goose migration tool if needed
go install github.com/pressly/goose/v3/cmd/goose@latest

# Apply migrations
goose -dir sql/schema postgres "postgres://localhost:5432/rss?sslmode=disable" up
```

## ⚙️ Configuration

Gator uses a configuration file located at `~/.gatorconfig.json`:

```json
{
  "current_user_name": "",
  "database_url": "postgres://localhost:5432/rss?sslmode=disable"
}
```

You can create this file manually or let Gator create it with default values on first run.

## Commands Reference

### User Management

| Command | Description | Example |
|---------|-------------|---------|
| `register` | Create a new user account | `gator register username` |
| `login` | Switch to another user | `gator login username` |

### Feed Management

| Command | Description | Example |
|---------|-------------|---------|
| `addfeed` | Add a new RSS feed | `gator addfeed "HackerNews" "https://news.ycombinator.com/rss"` |
| `feeds` | List all available feeds | `gator feeds` |
| `follow` | Follow an existing feed | `gator follow "https://news.ycombinator.com/rss"` |
| `following` | List feeds you're following | `gator following` |
| `unfollow` | Unfollow a feed | `gator unfollow "https://news.ycombinator.com/rss"` |

### Content Aggregation

| Command | Description | Example |
|---------|-------------|---------|
| `agg` | Start the aggregator (with interval) | `gator agg 5m` |
| `browse` | View posts from followed feeds | `gator browse` |
| `browse <limit>` | View specific number of posts | `gator browse 5` |

## Workflow

```
┌─────────────┐     ┌────────────┐     ┌──────────────┐     ┌──────────────┐     ┌────────────┐
│ 1. Register │ ──▶ │ 2. Add     │ ──▶ │ 3. Follow    │ ──▶ │ 4. Aggregate │ ──▶ │ 5. Browse  │
│    User     │     │    Feeds   │     │    Feeds     │     │    Posts     │     │    Posts   │
└─────────────┘     └────────────┘     └──────────────┘     └──────────────┘     └────────────┘
```

The aggregator runs in the background and continuously collects new posts from your followed feeds at the interval you specify.

## Usage Examples

### Setting Up and Following Feeds

```bash
# Create a user
gator register johndoe

# Add some popular tech feeds
gator addfeed "HackerNews" "https://news.ycombinator.com/rss"
gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
gator addfeed "Lobsters" "https://lobste.rs/rss"

# Start the aggregator in the background (updating every 10 minutes)
gator agg 10m &

# In another terminal, browse the latest posts
gator browse
```

## Development

For development, use:

```bash
go run . <command>
```

For production use, install the binary and use:

```bash
gator <command>
```

