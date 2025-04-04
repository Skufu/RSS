package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Skufu/RSS/internal/config"
)

// state holds application state that can be passed to command handlers
type state struct {
	config *config.Config
}

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Initialize application state
	s := &state{
		config: &cfg,
	}

	// Initialize commands
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register command handlers
	cmds.register("login", handlerLogin)

	// Process command line arguments
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		fmt.Println("Usage: gator <command> [args...]")
		fmt.Println("Available commands: login")
		os.Exit(1)
	}

	// Parse the command (skip first arg which is the program name)
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	// Run the command
	if err := cmds.run(s, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
