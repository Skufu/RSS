package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Skufu/RSS/internal/app"
	"github.com/Skufu/RSS/internal/database"
	"github.com/google/uuid"
)

// HandlerRegister handles the register command which creates a new user
func HandlerRegister(s *app.State, cmd app.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("register command requires a username argument")
	}

	username := cmd.Args[0]

	// Check if user already exists in database
	ctx := context.Background()
	_, err := s.Db.GetUserByName(ctx, username)
	if err == nil {
		// User already exists
		fmt.Printf("Error: User '%s' already exists\n", username)
		os.Exit(1)
	}

	// Create new user
	now := time.Now()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	}

	newUser, err := s.Db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set current user in config
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	// Print success message and user data
	fmt.Printf("User '%s' created successfully\n", username)
	fmt.Printf("User data: ID=%s, Created=%s, Name=%s\n",
		newUser.ID, newUser.CreatedAt.Format(time.RFC3339), newUser.Name)

	return nil
}
