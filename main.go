package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	client := &http.Client{}

	sessionState := app.NewState(cfg, dbQueries, client)

	cmd, err := commands.ParseArgs(os.Args...)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	err = commands.Run(&sessionState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
