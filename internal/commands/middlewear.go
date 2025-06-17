package commands

import (
	"context"
	"fmt"

	"github.com/viyan-md/gator_rss/internal/app"
	"github.com/viyan-md/gator_rss/internal/database"
)

func LoggedIn(
	handler func(s *app.State, cmd Command, user database.User) error,
) func(s *app.State, cmd Command) error {
	return func(s *app.State, cmd Command) error {
		name := s.Config.CurrentUserName
		if name == "" {
			return fmt.Errorf("you must be logged in")
		}

		user, err := s.DBQueries.GetUser(context.Background(), name)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
