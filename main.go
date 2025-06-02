package main

import (
	"fmt"

	"github.com/Xeninon/Gator/internal/config"
)

func main() {
	configs, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	configs.SetUser("Xeninon")
	configs, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(configs)
}
