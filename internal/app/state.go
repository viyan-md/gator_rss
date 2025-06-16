package app

import (
	"net/http"

	"github.com/viyan-md/gator_rss/internal/config"
	"github.com/viyan-md/gator_rss/internal/database"
)

type State struct {
	Config    *config.Config
	DBQueries *database.Queries
	Client    *http.Client
}

func NewState(cfg config.Config, dbq *database.Queries, c *http.Client) State {
	return State{
		Config:    &cfg,
		DBQueries: dbq,
		Client:    c,
	}
}
