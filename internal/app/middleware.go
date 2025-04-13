package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Skufu/RSS/internal/database"
)

// HandlerFunc defines the standard handler function signature
type HandlerFunc func(s *State, cmd Command) error

// AuthHandlerFunc defines the handler function signature for routes requiring authentication
type AuthHandlerFunc func(s *State, cmd Command, user database.User) error

// returns a standard handler function that checks for a logged-in user first.
func MiddlewareLoggedIn(handler AuthHandlerFunc) HandlerFunc {
	return func(s *State, cmd Command) error {
		// Check if a user is set in the config
		if s.Cfg.CurrentUserName == "" {
			return errors.New("no user set, please login first")
		}

		// Get the current user from the database
		ctx := context.Background()
		user, err := s.Db.GetUserByName(ctx, s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Call the handler with the user provided
		return handler(s, cmd, user)
	}
}
