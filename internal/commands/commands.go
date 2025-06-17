package commands

import (
	"errors"
	"fmt"

	"github.com/viyan-md/gator_rss/internal/app"
)

func Run(s *app.State, cmd Command) error {
	handler, ok := getCommands()[cmd.Name]
	if !ok {
		return errors.New("error: invalid command")
	}

	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func getCommands() map[string]func(*app.State, Command) error {
	return map[string]func(*app.State, Command) error{
		"login":     HandlerLogin,
		"register":  HandlerRegister,
		"reset":     HandlerReset,
		"users":     HandleUsers,
		"agg":       HandleAgg,
		"addfeed":   LoggedIn(HandleAddFeed),
		"feeds":     HandleGetFeeds,
		"follow":    LoggedIn(HandleFollowFeed),
		"following": LoggedIn(HandleListFollowing),
	}
}

func ParseArgs(args ...string) (Command, error) {
	if len(args) < 2 {
		return Command{}, errors.New("invalid input")
	}

	return Command{
		Name: args[1],
		Args: args[2:],
	}, nil
}
