package clippingsfeed

import (
	"fmt"
	"os"
	"path/filepath"
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

// FilterValidMetadata removes metadata items without required Source or Title fields
func FilterValidMetadata(metadata []Metadata) []Metadata {
	var filtered []Metadata
	for _, meta := range metadata {
		if meta.Source != "" && meta.Title != "" {
			filtered = append(filtered, meta)
		}
	}
	return filtered
}

// SortMetadataByCreated sorts metadata by Created field in descending order (newest first)
func SortMetadataByCreated(metadata []Metadata) {
	sort.Slice(metadata, func(i, j int) bool {
		return metadata[i].Created.After(metadata[j].Created)
	})
}

// LimitMetadataItems limits the number of items if maxItems is specified and positive
func LimitMetadataItems(metadata []Metadata, maxItems int) []Metadata {
	if maxItems > 0 && len(metadata) > maxItems {
		return metadata[:maxItems]
	}
	return metadata
}

func GenerateFeed(metadata []Metadata, config FeedConfig) (*feeds.Feed, error) {
	// Process metadata: filter, sort, and limit
	filteredMetadata := FilterValidMetadata(metadata)
	SortMetadataByCreated(filteredMetadata)
	processedMetadata := LimitMetadataItems(filteredMetadata, config.MaxItems)

	feed := &feeds.Feed{
		Title:       config.Title,
		Link:        &feeds.Link{Href: config.Link},
		Description: config.Description,
		Author:      &feeds.Author{Name: config.Author},
		Created:     config.Created,
	}

	for _, meta := range processedMetadata {

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

func WriteFeedToFile(feed *feeds.Feed, filename string) error {
	// Determine format from file extension
	ext := strings.ToLower(filepath.Ext(filename))
	var format string
	switch ext {
	case ".rss":
		format = "rss"
	case ".atom":
		format = "atom"
	case ".json":
		format = "json"
	default:
		return fmt.Errorf("unsupported file extension: %s (supported: .rss, .atom, .json)", ext)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// If there was already an error, keep it; otherwise use close error
			if err == nil {
				err = fmt.Errorf("failed to close file %s: %w", filename, closeErr)
			}
		}
	}()

	switch format {
	case "rss":
		err = feed.WriteRss(file)
	case "atom":
		err = feed.WriteAtom(file)
	case "json":
		err = feed.WriteJSON(file)
	}

	if err != nil {
		return fmt.Errorf("failed to write %s feed to %s: %w", format, filename, err)
	}

	return nil
}
