package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/viyan-md/gator_rss/internal/app"
	"github.com/viyan-md/gator_rss/internal/commands"
	"github.com/viyan-md/gator_rss/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	sessionState := app.NewState(cfg)

	sessionCommandsList := commands.CommandsList{
		Commands: make(map[string]func(*app.State, commands.Command) error),
	}

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("Error: invalid input\n")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: args[0],
		Args: args[1:],
	}
	sessionCommandsList.Register("login", commands.HandlerLogin)

	err = sessionCommandsList.Run(&sessionState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
