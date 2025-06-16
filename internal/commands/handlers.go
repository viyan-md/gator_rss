package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/viyan-md/gator_rss/internal/app"
	"github.com/viyan-md/gator_rss/internal/database"
)

func HandlerLogin(s *app.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("login command requires username argument")
	}

	username := cmd.Args[0]

	user, err := s.DBQueries.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("Error: user doesn't exist.")
		os.Exit(1)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %v\n", user.Name)

	return nil
}

func HandlerRegister(s *app.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("register command requires username argument")
	}

	username := cmd.Args[0]
	user, err := s.DBQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: user already exists.")
			os.Exit(1)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	fmt.Printf("%s user has been set!\n", user.Name)
	return nil
}

func HandlerReset(s *app.State, cmd Command) error {
	err := s.DBQueries.ResetUsers(context.Background())

	if err != nil {
		fmt.Printf("failed to reset users table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Reset successful")
	return nil
}

func HandleUsers(s *app.State, cmd Command) error {
	users, err := s.DBQueries.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("failed to load users: %v", err)
		os.Exit(1)
	}

	if len(users) < 1 {
		fmt.Println("empty")
		return nil
	}

	for _, user := range users {
		fmt.Printf("* %s ", user.Name)
		if user.Name == s.Config.CurrentUserName {
			fmt.Print("(current)")
		}
		fmt.Println()
	}

	return nil
}
