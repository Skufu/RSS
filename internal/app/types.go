package app

import (
	"github.com/Skufu/RSS/internal/config"
	"github.com/Skufu/RSS/internal/database"
)

// State holds application state that can be passed to command handlers
type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

// Command represents a CLI command with its name and arguments
type Command struct {
	Name string
	Args []string
}
