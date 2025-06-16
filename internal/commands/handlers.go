package commands

import (
	"errors"
	"fmt"

	"github.com/viyan-md/gator_rss/internal/app"
)

func HandlerLogin(s *app.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("login command requires username argument")
	}

	username := cmd.Args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("%s user has been set!", username)
	return nil
}
