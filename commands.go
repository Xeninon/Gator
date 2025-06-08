package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Xeninon/Gator/internal/database"
	"github.com/google/uuid"

	"github.com/Xeninon/Gator/internal/config"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) {
	handler, ok := c.commandMap[cmd.name]
	if !ok {
		fmt.Println("command not registered")
		os.Exit(1)
	}

	if err := handler(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return errors.New("username is not registered")
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Println("user has been set to " + user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("name is required")
	}

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0],
		},
	)
	if err != nil {
		return errors.New("username is taken")
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("User was created with info:%v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.Current_user_name {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("duration string required")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return errors.New("requires valid duration string")
	}

	fmt.Println("Collecting feeds every " + timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return errors.New("name and url required")
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0],
			Url:       cmd.arguments[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return errors.New("url already in database")
	}

	s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)

	fmt.Printf("Feed was created with info:%v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("feed: %v, url: %v, creator: %v\n", feed.Name, feed.Url, feed.Creator)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return errors.New("url required")
	}

	feed, err := s.db.GetFeedForURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return errors.New("feed not in database")
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return errors.New("already following feed")
	}

	fmt.Println(feedFollow.UserName + "now following" + feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return errors.New("url required")
	}

	err := s.db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			Url:    cmd.arguments[0],
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		cmd.arguments[0] = "2"
	}

	limit, err := strconv.Atoi(cmd.arguments[0])
	if err != nil {
		return errors.New("provided argument isn't an integer")
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post.Title)
		fmt.Println("| " + post.Url)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return errors.New("account not registered in database")
		}

		if err := handler(s, cmd, user); err != nil {
			return err
		}

		return nil
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	s.db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID: feed.ID,
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
	)

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rss.Channel.Item {
		pubdate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				PublishedAt: pubdate,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				fmt.Println(err)
			}
		}
	}

	return nil
}
