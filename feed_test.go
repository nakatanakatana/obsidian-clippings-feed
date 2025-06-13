package clippingsfeed_test

import (
	"testing"
	"time"

	clippingsfeed "github.com/nakatanakatana/obsidian-clippings-feed"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestGenerateFeed(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	for name, tt := range map[string]struct {
		goldenFilename string
		metadata       []clippingsfeed.Metadata
		config         clippingsfeed.FeedConfig
	}{
		"single item": {
			goldenFilename: "feed_single",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "Test Article",
					Site:        "example.com",
					Source:      "https://example.com/article1",
					Author:      []string{"John Doe"},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "This is a test article",
					Tags:        []string{"go", "testing"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:       "Test Feed",
				Link:        "https://example.com/feed",
				Description: "Test RSS Feed",
				Author:      "Feed Author",
				Created:     baseTime,
			},
		},
		"multiple items": {
			goldenFilename: "feed_multiple",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "First Article",
					Site:        "example.com",
					Source:      "https://example.com/article1",
					Author:      []string{"Author One"},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "First test article",
					Tags:        []string{"go", "web"},
				},
				{
					Title:       "Second Article",
					Site:        "example.com",
					Source:      "https://example.com/article2",
					Author:      []string{"Author Two"},
					Published:   "2025-06-02",
					Created:     baseTime.Add(24 * time.Hour),
					Description: "Second test article",
					Tags:        []string{"testing", "feed"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:       "Multi Article Feed",
				Link:        "https://example.com/multi-feed",
				Description: "Feed with multiple articles",
				Author:      "Multi Feed Author",
				Created:     baseTime,
			},
		},
		"empty description": {
			goldenFilename: "feed_empty_desc",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "No Description Article",
					Site:        "example.com",
					Source:      "https://example.com/article",
					Author:      []string{},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "",
					Tags:        []string{"minimal"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:       "Minimal Feed",
				Link:        "https://example.com/minimal",
				Description: "Feed with minimal data",
				Author:      "Minimal Author",
				Created:     baseTime,
			},
		},
		"limited items": {
			goldenFilename: "feed_limited",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "First Article",
					Site:        "example.com",
					Source:      "https://example.com/article1",
					Author:      []string{"Author One"},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "First test article",
					Tags:        []string{"go", "web"},
				},
				{
					Title:       "Second Article",
					Site:        "example.com",
					Source:      "https://example.com/article2",
					Author:      []string{"Author Two"},
					Published:   "2025-06-02",
					Created:     baseTime.Add(24 * time.Hour),
					Description: "Second test article",
					Tags:        []string{"testing", "feed"},
				},
				{
					Title:       "Third Article",
					Site:        "example.com",
					Source:      "https://example.com/article3",
					Author:      []string{"Author Three"},
					Published:   "2025-06-03",
					Created:     baseTime.Add(48 * time.Hour),
					Description: "Third test article",
					Tags:        []string{"limit", "test"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:       "Limited Feed",
				Link:        "https://example.com/limited-feed",
				Description: "Feed with item limit",
				Author:      "Limited Feed Author",
				Created:     baseTime,
				MaxItems:    2,
			},
		},
		"hidden description": {
			goldenFilename: "feed_hidden_desc",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "Test Article",
					Site:        "example.com",
					Source:      "https://example.com/article1",
					Author:      []string{"John Doe"},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "This description should be hidden",
					Tags:        []string{"go", "testing"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:           "Hidden Desc Feed",
				Link:            "https://example.com/hidden-desc",
				Description:     "Feed with hidden descriptions",
				Author:          "Hidden Desc Author",
				Created:         baseTime,
				HideDescription: true,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			feed, err := clippingsfeed.GenerateFeed(tt.metadata, tt.config)
			assert.NilError(t, err)

			result, err := feed.ToAtom()
			assert.NilError(t, err)

			golden.Assert(t, result, tt.goldenFilename)
		})
	}
}
