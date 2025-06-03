package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Xeninon/Gator/internal/database"

	"github.com/Xeninon/Gator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	currentState := state{
		dbQueries,
		&cfg,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	if len(os.Args) < 2 {
		fmt.Println("no command given")
		os.Exit(1)
	}

	args := make([]string, 0)
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmd := command{
		os.Args[1],
		args,
	}

	cmds.run(&currentState, cmd)
}
