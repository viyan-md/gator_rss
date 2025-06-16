package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/viyan-md/gator_rss/internal/app"
	"github.com/viyan-md/gator_rss/internal/commands"
	"github.com/viyan-md/gator_rss/internal/config"
	"github.com/viyan-md/gator_rss/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	sessionState := app.NewState(cfg, dbQueries)

	sessionCommandsList := commands.CommandsList{
		Commands: make(map[string]func(*app.State, commands.Command) error),
	}

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Printf("Error: invalid input\n")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: args[0],
		Args: args[1:],
	}

	sessionCommandsList.Register("login", commands.HandlerLogin)
	sessionCommandsList.Register("register", commands.HandlerRegister)
	sessionCommandsList.Register("reset", commands.HandlerReset)
	sessionCommandsList.Register("users", commands.HandleUsers)

	err = sessionCommandsList.Run(&sessionState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
