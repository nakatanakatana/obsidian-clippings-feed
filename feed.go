package clippingsfeed

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type FeedConfig struct {
	Title       string
	Link        string
	Description string
	Author      string
	Created     time.Time
}

func GenerateFeed(metadata []Metadata, config FeedConfig) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title:       config.Title,
		Link:        &feeds.Link{Href: config.Link},
		Description: config.Description,
		Author:      &feeds.Author{Name: config.Author},
		Created:     config.Created,
	}

	for _, meta := range metadata {
		var authorName string
		if len(meta.Author) > 0 {
			authorName = strings.Join(meta.Author, ", ")
		}

		item := &feeds.Item{
			Title:       meta.Title,
			Link:        &feeds.Link{Href: meta.Source},
			Description: meta.Description,
			Author:      &feeds.Author{Name: authorName},
			Created:     meta.Created,
			Id:          meta.Source,
		}

		if len(meta.Author) > 0 {
			item.Description += fmt.Sprintf("\n\nAuthor(s): %s", strings.Join(meta.Author, ", "))
		}

		if len(meta.Tags) > 0 {
			item.Description += fmt.Sprintf("\n\nTags: %v", meta.Tags)
		}

		if meta.Site != "" {
			item.Description += fmt.Sprintf("\n\nSite: %s", meta.Site)
		}

		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

func WriteFeedToFile(feed *feeds.Feed, filename string, format string) error {
	var feedContent string
	var err error

	switch format {
	case "rss":
		feedContent, err = feed.ToRss()
	case "atom":
		feedContent, err = feed.ToAtom()
	case "json":
		feedContent, err = feed.ToJSON()
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to generate feed: %w", err)
	}

	err = os.WriteFile(filename, []byte(feedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
