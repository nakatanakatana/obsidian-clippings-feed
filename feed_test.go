package clippingsfeed_test

import (
	"os"
	"path/filepath"
	"strings"
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
		"filtered items": {
			goldenFilename: "feed_filtered",
			metadata: []clippingsfeed.Metadata{
				{
					Title:       "Valid Article",
					Site:        "example.com",
					Source:      "https://example.com/valid",
					Author:      []string{"Valid Author"},
					Published:   "2025-06-01",
					Created:     baseTime,
					Description: "This article has both title and source",
					Tags:        []string{"valid"},
				},
				{
					Title:       "", // Empty title - should be filtered
					Site:        "example.com",
					Source:      "https://example.com/no-title",
					Author:      []string{"No Title Author"},
					Published:   "2025-06-02",
					Created:     baseTime.Add(24 * time.Hour),
					Description: "This article has no title",
					Tags:        []string{"invalid"},
				},
				{
					Title:       "No Source Article",
					Site:        "example.com",
					Source:      "", // Empty source - should be filtered
					Author:      []string{"No Source Author"},
					Published:   "2025-06-03",
					Created:     baseTime.Add(48 * time.Hour),
					Description: "This article has no source",
					Tags:        []string{"invalid"},
				},
			},
			config: clippingsfeed.FeedConfig{
				Title:       "Filtered Feed",
				Link:        "https://example.com/filtered",
				Description: "Feed with filtered items",
				Author:      "Filtered Author",
				Created:     baseTime,
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

func TestFilterValidMetadata(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	input := []clippingsfeed.Metadata{
		{
			Title:   "Valid Article",
			Source:  "https://example.com/valid",
			Created: baseTime,
		},
		{
			Title:   "", // Empty title - should be filtered
			Source:  "https://example.com/no-title",
			Created: baseTime,
		},
		{
			Title:   "No Source Article",
			Source:  "", // Empty source - should be filtered
			Created: baseTime,
		},
		{
			Title:   "Another Valid Article",
			Source:  "https://example.com/valid2",
			Created: baseTime,
		},
	}

	result := clippingsfeed.FilterValidMetadata(input)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Valid Article", result[0].Title)
	assert.Equal(t, "Another Valid Article", result[1].Title)
}

func TestSortMetadataByCreated(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	metadata := []clippingsfeed.Metadata{
		{
			Title:   "Oldest",
			Source:  "https://example.com/1",
			Created: baseTime,
		},
		{
			Title:   "Newest",
			Source:  "https://example.com/3",
			Created: baseTime.Add(48 * time.Hour),
		},
		{
			Title:   "Middle",
			Source:  "https://example.com/2",
			Created: baseTime.Add(24 * time.Hour),
		},
	}

	clippingsfeed.SortMetadataByCreated(metadata)

	assert.Equal(t, "Newest", metadata[0].Title)
	assert.Equal(t, "Middle", metadata[1].Title)
	assert.Equal(t, "Oldest", metadata[2].Title)
}

func TestLimitMetadataItems(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	metadata := []clippingsfeed.Metadata{
		{Title: "Item 1", Source: "https://example.com/1", Created: baseTime},
		{Title: "Item 2", Source: "https://example.com/2", Created: baseTime},
		{Title: "Item 3", Source: "https://example.com/3", Created: baseTime},
		{Title: "Item 4", Source: "https://example.com/4", Created: baseTime},
	}

	// Test with limit
	result := clippingsfeed.LimitMetadataItems(metadata, 2)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Item 1", result[0].Title)
	assert.Equal(t, "Item 2", result[1].Title)

	// Test without limit (0)
	result = clippingsfeed.LimitMetadataItems(metadata, 0)
	assert.Equal(t, 4, len(result))

	// Test with limit larger than input
	result = clippingsfeed.LimitMetadataItems(metadata, 10)
	assert.Equal(t, 4, len(result))
}

func TestWriteFeedToFile(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	feed, err := clippingsfeed.GenerateFeed([]clippingsfeed.Metadata{
		{
			Title:       "Test Article",
			Site:        "example.com",
			Source:      "https://example.com/test",
			Author:      []string{"Test Author"},
			Published:   "2025-06-01",
			Created:     baseTime,
			Description: "Test description",
			Tags:        []string{"test"},
		},
	}, clippingsfeed.FeedConfig{
		Title:       "Test Feed",
		Link:        "https://example.com/feed",
		Description: "Test Description",
		Author:      "Test Author",
		Created:     baseTime,
	})
	assert.NilError(t, err)

	tests := []struct {
		name     string
		filename string
		contains string
	}{
		{
			name:     "RSS format",
			filename: "test.rss",
			contains: "<?xml version=\"1.0\" encoding=\"UTF-8\"?><rss",
		},
		{
			name:     "Atom format",
			filename: "test.atom",
			contains: "<?xml version=\"1.0\" encoding=\"UTF-8\"?><feed",
		},
		{
			name:     "JSON format",
			filename: "test.json",
			contains: `"version":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tt.filename)

			err := clippingsfeed.WriteFeedToFile(feed, tmpFile)
			assert.NilError(t, err)

			// Verify file exists and contains expected content
			content, err := os.ReadFile(tmpFile)
			assert.NilError(t, err)
			assert.Assert(t, strings.Contains(string(content), tt.contains))
			assert.Assert(t, strings.Contains(string(content), "Test Feed"))
		})
	}
}

func TestWriteFeedToFileUnsupportedExtension(t *testing.T) {
	baseTime := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	feed, err := clippingsfeed.GenerateFeed([]clippingsfeed.Metadata{
		{
			Title:   "Test Article",
			Source:  "https://example.com/test",
			Created: baseTime,
		},
	}, clippingsfeed.FeedConfig{
		Title:   "Test Feed",
		Link:    "https://example.com/feed",
		Created: baseTime,
	})
	assert.NilError(t, err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.xml")

	err = clippingsfeed.WriteFeedToFile(feed, tmpFile)
	assert.ErrorContains(t, err, "unsupported file extension: .xml")
}
