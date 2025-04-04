package main

import (
	"errors"
	"fmt"
)

// handlerLogin handles the login command which sets the current user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login command requires a username argument")
	}

	username := cmd.args[0]
	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User set to '%s'\n", username)
	return nil
}
