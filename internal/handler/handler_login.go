package handler

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Skufu/RSS/internal/app"
)

// HandlerLogin handles the login command which sets the current user
func HandlerLogin(s *app.State, cmd app.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login command requires a username argument")
	}

	username := cmd.Args[0]

	// Check if user exists in database
	ctx := context.Background()
	_, err := s.Db.GetUserByName(ctx, username)
	if err != nil {
		// User doesn't exist
		fmt.Printf("Error: User '%s' doesn't exist\n", username)
		os.Exit(1)
	}

	// Update config with current user
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User set to '%s'\n", username)
	return nil
}
