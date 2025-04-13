package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/Skufu/RSS/internal/app"
)

// HandlerUsers handles the users command which lists all users
func HandlerUsers(s *app.State, cmd app.Command) error {
	// Get all users from the database
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
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
	currentUser := s.Cfg.CurrentUserName

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
