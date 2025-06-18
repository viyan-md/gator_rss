package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/viyan-md/gator_rss/internal/app"
	"github.com/viyan-md/gator_rss/internal/database"
	"github.com/viyan-md/gator_rss/internal/rss"
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

func HandleAgg(s *app.State, cmd Command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		rss.ScrapeFeeds(s)
	}
}

func HandleAddFeed(s *app.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("addfeed command requires name and url arguments")
	}

	newFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	newFeed, err := s.DBQueries.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}

	_, err = s.DBQueries.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	printFeed(&newFeed)
	return nil
}

func printFeed(f *database.Feed) {
	fmt.Printf("ID: %v\n", f.ID)
	fmt.Printf("CreatedAt: %v\n", f.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", f.UpdatedAt)
	fmt.Printf("Name: %v\n", f.Name)
	fmt.Printf("Url: %v\n", f.Url)
	fmt.Printf("UserID: %v\n", f.UserID)
}

func HandleGetFeeds(s *app.State, cmd Command) error {
	feedsList, err := s.DBQueries.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	printFeedList(feedsList)
	return nil
}

func printFeedList(fl []database.GetFeedsRow) {
	for _, row := range fl {
		fmt.Printf("Name:     %v\n", row.FeedName)
		fmt.Printf("URL:      %v\n", row.FeedUrl)
		fmt.Printf("User:     %v\n", row.UserName)
	}
}

func HandleFollowFeed(s *app.State, cmd Command, user database.User) error {
	feed, err := s.DBQueries.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	createdRow, err := s.DBQueries.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	printRow(createdRow.FeedName, createdRow.UserName)
	return nil
}

func HandleListFollowing(s *app.State, cmd Command, user database.User) error {
	followList, err := s.DBQueries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, fname := range followList {
		fmt.Printf(" - %v\n", fname)
	}

	return nil
}

func printRow(args ...string) {
	for _, val := range args {
		fmt.Println(val)
	}
}

func HandleUnfollowFeed(s *app.State, cmd Command, user database.User) error {
	unfollowParams := database.UnfollowFeedParams{
		ID:  user.ID,
		Url: cmd.Args[0],
	}

	unfollowed, err := s.DBQueries.UnfollowFeed(context.Background(), unfollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed: \n%v\n", unfollowed)
	return nil
}
