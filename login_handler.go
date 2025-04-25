package main

import (
	"fmt"
	"os"
)

func handerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	user := cmd.arguments[0]

	s.config.CurrentUserName = user

	fmt.Printf("user set as %v\n", user)

	return nil
}
