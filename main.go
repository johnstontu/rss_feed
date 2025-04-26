package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/johnstontu/rss_feed/internal/config"
	"github.com/johnstontu/rss_feed/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}

	var state State
	state.config = &cfg

	dbURL := state.config.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("%v", err)
	}

	dbQueries := database.New(db)
	state.db = dbQueries

	cmds := Commands{
		command: make(map[string]func(*State, Command) error),
	}

	cmds.register("login", handerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	input := os.Args
	if len(input) < 2 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}
	name := input[1]
	args := input[2:]

	command := Command{
		name:      name,
		arguments: args,
	}

	cmds.run(&state, command)

	fmt.Printf("%+v\n", cfg)

	config.Write(cfg)

}
