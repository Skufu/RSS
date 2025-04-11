package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Skufu/RSS/internal/config"
	"github.com/Skufu/RSS/internal/database"
	_ "github.com/lib/pq"
)

// state holds application state that can be passed to command handlers
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Load Database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create database queries
	dbQueries := database.New(db)

	// Initialize application state
	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Initialize commands
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register command handlers
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	// Process command line arguments
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		fmt.Println("Usage: gator <command> [args...]")
		fmt.Println("Available commands: login, register, reset, users, agg, addfeed, feeds")
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
