package main

import (
	"fmt"
)

// command represents a CLI command with its name and arguments
type command struct {
	name string
	args []string
}

// commands holds all the commands the CLI can handle
type commands struct {
	handlers map[string]func(*state, command) error
}

// register adds a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// run executes a given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}
