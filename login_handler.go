package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/johnstontu/rss_feed/internal/database"
)

func handerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	user := cmd.arguments[0]

	_, err := s.db.GetUser(
		context.Background(),
		user,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "User doesn't exist %v\n", err)
		os.Exit(1)
	}

	s.config.CurrentUserName = user

	fmt.Printf("user set as %v\n", user)

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	name := cmd.arguments[0]

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user %v\n", err)
		os.Exit(1)
	}

	s.config.CurrentUserName = user.Name

	return nil
}
