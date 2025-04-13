package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/Skufu/RSS/internal/app"
)

// HandlerReset handles the reset command which deletes all users
func HandlerReset(s *app.State, cmd app.Command) error {
	// Delete all users from the database
	ctx := context.Background()
	err := s.Db.ResetUsers(ctx)
	if err != nil {
		fmt.Printf("Error: Failed to reset database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database reset successful. All users have been deleted.")
	return nil
}
