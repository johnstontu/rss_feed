package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/johnstontu/rss_feed/internal/database"
)

func handerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	user := cmd.arguments[0]

	_, err := s.db.GetUser(
		context.Background(),
		user,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "User doesn't exist %v\n", err)
		os.Exit(1)
	}

	s.config.CurrentUserName = user

	fmt.Printf("user set as %v\n", user)

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	name := cmd.arguments[0]

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user %v\n", err)
		os.Exit(1)
	}

	s.config.CurrentUserName = user.Name

	return nil
}

func handlerReset(s *State, cmd Command) error {
	s.db.DeleteUsers(
		context.Background(),
	)

	return nil
}

func handlerUsers(s *State, cmd Command) error {

	users, err := s.db.GetUsers(
		context.Background(),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users %v\n", err)
		os.Exit(1)
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("%+v (current)", user.Name)
		} else {
			fmt.Println(user.Name)
		}

	}

	return nil

}

func handlerAgg(s *State, cmd Command) error {

	if len(cmd.arguments) < 1 || len(cmd.arguments) > 2 {
		fmt.Println("Wrong number of arguments")
		os.Exit(1)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Println("Invalid duration")
		os.Exit(1)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *State) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func handlerAddFeed(s *State, cmd Command) error {

	if len(cmd.arguments) < 2 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	name := cmd.arguments[0]
	url := cmd.arguments[1]

	currentUser, err := s.db.GetUser(
		context.Background(),
		s.config.CurrentUserName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current user %v\n", err)
		os.Exit(1)
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    currentUser.ID,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating feed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", feed)

	handlerFollow(s, cmd)

	return nil
}

func handlerFeeds(s *State, cmd Command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current user %v\n", err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Printf("%v\n", feed.FeedName)
		fmt.Printf("%v\n", feed.Url)
		fmt.Printf("%v\n", feed.UserName.String)
	}

	return nil
}

func handlerFollow(s *State, cmd Command) error {

	if len(cmd.arguments) < 1 {
		fmt.Println("Needs more input arguments")
		os.Exit(1)
	}

	var url string

	if len(cmd.arguments) == 2 {
		url = cmd.arguments[1]
	} else {
		url = cmd.arguments[0]
	}

	currentUser, err := s.db.GetUser(
		context.Background(),
		s.config.CurrentUserName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current user %v\n", err)
		os.Exit(1)
	}

	feed, err := s.db.GetFeed(
		context.Background(),
		url,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching feed %v\n", err)
		os.Exit(1)
	}

	s.db.CreateFeedFollows(
		context.Background(),
		database.CreateFeedFollowsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    currentUser.ID,
			FeedID:    feed.FeedID,
		},
	)

	return nil

}

func handlerFollowing(s *State, cmd Command) error {

	currentUser, err := s.db.GetUser(
		context.Background(),
		s.config.CurrentUserName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current user %v\n", err)
		os.Exit(1)
	}

	following, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		currentUser.Name,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching follows for user %v\n", err)
		os.Exit(1)
	}

	for _, follow := range following {
		fmt.Printf("%v\n", follow.UserName)
		fmt.Printf("%v\n", follow.FeedName)

	}

	return nil

}

func handlerUnfollow(s *State, cmd Command) error {

	currentUser, err := s.db.GetUser(
		context.Background(),
		s.config.CurrentUserName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching current user %v\n", err)
		os.Exit(1)
	}

	if len(cmd.arguments) != 1 {
		fmt.Println("Only takes one URL argument")
		os.Exit(1)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("Error getting feed")
		os.Exit(1)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feed.FeedID,
	})
	if err != nil {
		fmt.Println("Error deleting feed")
		os.Exit(1)
	}

	fmt.Printf("%s unfollowed successfully\n", feed.FeedName)
	return nil

}
