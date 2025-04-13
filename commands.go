package main

import (
	"fmt"

	"github.com/Skufu/RSS/internal/app"
)

// command represents a CLI command with its name and arguments
// type command struct {
// 	name string
// 	args []string
// }

// Commands holds all the commands the CLI can handle
type Commands struct {
	handlers map[string]func(*app.State, app.Command) error
}

// NewCommands creates a new Commands instance
func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*app.State, app.Command) error),
	}
}

// Register adds a new handler function for a command name
func (c *Commands) Register(name string, f func(*app.State, app.Command) error) {
	c.handlers[name] = f
}

// Run executes a given command with the provided state if it exists
func (c *Commands) Run(s *app.State, cmd app.Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}
