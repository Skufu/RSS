package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/config"
	"github.com/Skufu/RSS/internal/database"
	"github.com/Skufu/RSS/internal/handler"
	_ "github.com/lib/pq"
)

// State holds application state that can be passed to command handlers
// type state struct {
// 	db  *database.Queries
// 	cfg *config.Config
// }

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
	s := &app.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	// Initialize commands
	cmds := NewCommands()

	// Register command handlers
	// TODO: Update these handlers once they are moved
	cmds.Register("login", handler.HandlerLogin)
	cmds.Register("register", handler.HandlerRegister)
	cmds.Register("reset", handler.HandlerReset)
	cmds.Register("users", handler.HandlerUsers)
	cmds.Register("agg", handler.HandlerAgg)
	cmds.Register("addfeed", app.MiddlewareLoggedIn(handler.HandlerAddFeed))
	cmds.Register("feeds", handler.HandlerFeeds)
	cmds.Register("follow", app.MiddlewareLoggedIn(handler.HandlerFollow))
	cmds.Register("unfollow", app.MiddlewareLoggedIn(handler.HandlerUnfollow))
	cmds.Register("following", app.MiddlewareLoggedIn(handler.HandlerFollowing))
	cmds.Register("browse", app.MiddlewareLoggedIn(handler.HandlerBrowse))

	// Process command line arguments
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		fmt.Println("Usage: gator <command> [args...]")
		fmt.Println("Available commands: login, register, reset, users, agg, addfeed, feeds, follow, unfollow, following, browse")
		os.Exit(1)
	}

	// Parse the command (skip first arg which is the program name)
	cmd := app.Command{
		Name: args[1],
		Args: args[2:],
	}

	// Run the command
	if err := cmds.Run(s, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
