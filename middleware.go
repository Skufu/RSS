package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Skufu/RSS/internal/database"
)

// middlewareLoggedIn is a higher-order function that takes a handler requiring a logged-in user
// and returns a standard handler function that checks for a logged-in user first.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Check if a user is set in the config
		if s.cfg.CurrentUserName == "" {
			return errors.New("no user set, please login first")
		}

		// Get the current user from the database
		ctx := context.Background()
		user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Call the handler with the user provided
		return handler(s, cmd, user)
	}
}
