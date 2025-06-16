package app

import "github.com/viyan-md/gator_rss/internal/config"

type State struct {
	Config *config.Config
}

func NewState(cfg config.Config) State {
	return State{
		Config: &cfg,
	}
}
