package app

import (
	"github.com/viyan-md/gator_rss/internal/config"
	"github.com/viyan-md/gator_rss/internal/database"
)

type State struct {
	Config    *config.Config
	DBQueries *database.Queries
}

func NewState(cfg config.Config, dbq *database.Queries) State {
	return State{
		Config:    &cfg,
		DBQueries: dbq,
	}
}
