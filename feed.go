package clippingsfeed

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type FeedConfig struct {
	Title           string
	Link            string
	Description     string
	Author          string
	Created         time.Time
	MaxItems        int
	HideDescription bool
}

func GenerateFeed(metadata []Metadata, config FeedConfig) (*feeds.Feed, error) {
	// Sort metadata by Created field in descending order (newest first)
	sort.Slice(metadata, func(i, j int) bool {
		return metadata[i].Created.After(metadata[j].Created)
	})

	// Limit the number of items if MaxItems is specified and positive
	if config.MaxItems > 0 && len(metadata) > config.MaxItems {
		metadata = metadata[:config.MaxItems]
	}

	feed := &feeds.Feed{
		Title:       config.Title,
		Link:        &feeds.Link{Href: config.Link},
		Description: config.Description,
		Author:      &feeds.Author{Name: config.Author},
		Created:     config.Created,
	}

	for _, meta := range metadata {
		// Skip items without required Source or Title
		if meta.Source == "" || meta.Title == "" {
			continue
		}

		var authorName string
		if len(meta.Author) > 0 {
			authorName = strings.Join(meta.Author, ", ")
		}

		var description string
		if !config.HideDescription {
			description = meta.Description

			if len(meta.Author) > 0 {
				description += fmt.Sprintf("\n\nAuthor(s): %s", strings.Join(meta.Author, ", "))
			}

			if len(meta.Tags) > 0 {
				description += fmt.Sprintf("\n\nTags: %v", meta.Tags)
			}

			if meta.Site != "" {
				description += fmt.Sprintf("\n\nSite: %s", meta.Site)
			}
		}

		item := &feeds.Item{
			Title:       meta.Title,
			Link:        &feeds.Link{Href: meta.Source},
			Description: description,
			Author:      &feeds.Author{Name: authorName},
			Created:     meta.Created,
			Id:          meta.Source,
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
