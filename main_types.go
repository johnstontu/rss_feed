package main

import (
	"github.com/johnstontu/rss_feed/internal/config"
	"github.com/johnstontu/rss_feed/internal/database"
)

type State struct {
	db     *database.Queries
	config *config.Config
}

type Command struct {
	name      string
	arguments []string
}

type Commands struct {
	command map[string]func(*State, Command) error
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.command[name] = f
}

func (c *Commands) run(s *State, cmd Command) error {
	function := c.command[cmd.name]
	function(s, cmd)
	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
