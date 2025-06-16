package commands

import "github.com/viyan-md/gator_rss/internal/app"

type Command struct {
	Name string
	Args []string
}

type CommandsList struct {
	Commands map[string]func(*app.State, Command) error
}
