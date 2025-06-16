package commands

import (
	"errors"
	"fmt"

	"github.com/viyan-md/gator_rss/internal/app"
)

func (c *CommandsList) Run(s *app.State, cmd Command) error {
	handler, ok := c.Commands[cmd.Name]
	if !ok {
		return errors.New("invalid command")
	}

	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	return nil
}

func (c *CommandsList) Register(name string, f func(*app.State, Command) error) {
	c.Commands[name] = f
}
