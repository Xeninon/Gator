package main

import (
	"fmt"
	"os"

	"github.com/Xeninon/Gator/internal/config"
)

func main() {
	configs, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	currentState := state{
		&configs,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
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
