package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Xeninon/Gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) {
	handler, ok := c.commandMap[cmd.name]
	if !ok {
		fmt.Println("command not registered")
		os.Exit(1)
	}

	if err := handler(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}

	if err := s.config.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	fmt.Println("user has been set to " + s.config.Current_user_name)
	return nil
}
