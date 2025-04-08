package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Skufu/RSS/internal/database"
	"github.com/google/uuid"
)

// handlerLogin handles the login command which sets the current user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login command requires a username argument")
	}

	username := cmd.args[0]

	// Check if user exists in database
	ctx := context.Background()
	_, err := s.db.GetUserByName(ctx, username)
	if err != nil {
		// User doesn't exist
		fmt.Printf("Error: User '%s' doesn't exist\n", username)
		os.Exit(1)
	}

	// Update config with current user
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User set to '%s'\n", username)
	return nil
}

// handlerRegister handles the register command which creates a new user
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register command requires a username argument")
	}

	username := cmd.args[0]

	// Check if user already exists in database
	ctx := context.Background()
	_, err := s.db.GetUserByName(ctx, username)
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

	newUser, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set current user in config
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	// Print success message and user data
	fmt.Printf("User '%s' created successfully\n", username)
	fmt.Printf("User data: ID=%s, Created=%s, Name=%s\n",
		newUser.ID, newUser.CreatedAt.Format(time.RFC3339), newUser.Name)

	return nil
}

// handlerReset handles the reset command which deletes all users
func handlerReset(s *state, cmd command) error {
	// Delete all users from the database
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		fmt.Printf("Error: Failed to reset database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database reset successful. All users have been deleted.")
	return nil
}

// handlerUsers handles the users command which lists all users
func handlerUsers(s *state, cmd command) error {
	// Get all users from the database
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Printf("Error: Failed to get users: %v\n", err)
		os.Exit(1)
	}

	// If no users, print a message
	if len(users) == 0 {
		fmt.Println("No users found in the database.")
		return nil
	}

	// Get the current user from config
	currentUser := s.cfg.CurrentUserName

	// Print each user, marking the current one
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
