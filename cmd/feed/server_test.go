package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	clippingsfeed "github.com/nakatanakatana/obsidian-clippings-feed"
	"gotest.tools/v3/golden"
)

func TestGenerateIndexHTML(t *testing.T) {
	tests := []struct {
		name     string
		testData string
		config   Config
	}{
		{
			name:     "single_item",
			testData: "success",
			config: Config{
				FeedTitle:     "Test Feed",
				FeedDesc:      "Test Description",
				FeedLink:      "http://example.com",
				FeedAuthor:    "Test Author",
				TargetDir:     "/test/dir",
				MaxItems:      50,
				DebounceDelay: 10 * time.Second,
			},
		},
		{
			name:     "multiple_items",
			testData: "feed_multiple",
			config: Config{
				FeedTitle:     "Multi Item Feed",
				FeedDesc:      "Feed with multiple items",
				FeedLink:      "http://example.com",
				FeedAuthor:    "Test Author",
				TargetDir:     "/test/dir",
				MaxItems:      50,
				DebounceDelay: 10 * time.Second,
			},
		},
		{
			name:     "empty_description",
			testData: "feed_empty_desc",
			config: Config{
				FeedTitle:     "Empty Desc Feed",
				FeedDesc:      "Feed with empty description items",
				FeedLink:      "http://example.com",
				FeedAuthor:    "Test Author",
				TargetDir:     "/test/dir",
				MaxItems:      50,
				DebounceDelay: 10 * time.Second,
			},
		},
		{
			name:     "hidden_description",
			testData: "success",
			config: Config{
				FeedTitle:       "Hidden Desc Feed",
				FeedDesc:        "Feed with hidden descriptions",
				FeedLink:        "http://example.com",
				FeedAuthor:      "Test Author",
				TargetDir:       "/test/dir",
				MaxItems:        50,
				DebounceDelay:   10 * time.Second,
				HideDescription: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tmpDir := t.TempDir()

			// Create generator with test config
			generator := &FeedGenerator{
				config: tt.config,
				tmpDir: tmpDir,
				parser: clippingsfeed.CreateParser(),
			}

			// Load test data
			metadata, err := loadTestData(t, tt.testData)
			if err != nil {
				t.Fatalf("Failed to load test data: %v", err)
			}

			// Generate HTML using the helper method
			outputFile := filepath.Join(tmpDir, "index.html")
			err = generator.generateIndexHTMLFromMetadata(outputFile, metadata)
			if err != nil {
				t.Fatalf("GenerateIndexHTML failed: %v", err)
			}

			// Read generated HTML
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read generated HTML: %v", err)
			}

			// Normalize timestamp for consistent testing
			normalizedContent := normalizeTimestamp(string(content))

			// Compare with golden file
			golden.Assert(t, normalizedContent, "index_"+tt.name+".html")
		})
	}
}

// loadTestData loads test metadata from testdata files
func loadTestData(t *testing.T, testDataName string) ([]clippingsfeed.Metadata, error) {
	t.Helper()

	switch testDataName {
	case "success":
		return []clippingsfeed.Metadata{
			{
				Title:       "Test Article",
				Site:        "example.com",
				Source:      "https://example.com/article",
				Author:      []string{"John Doe"},
				Published:   "2023-01-01",
				Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "This is a test article description.",
				Tags:        []string{"test", "article"},
			},
		}, nil
	case "feed_multiple":
		return []clippingsfeed.Metadata{
			{
				Title:       "First Article",
				Site:        "example.com",
				Source:      "https://example.com/first",
				Author:      []string{"John Doe", "Jane Smith"},
				Published:   "2023-01-01",
				Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "First test article.",
				Tags:        []string{"test", "first"},
			},
			{
Title:       "Second Article",
				Site:        "example.org",
				Source:      "https://example.org/second",
				Author:      []string{"Bob Wilson"},
				Published:   "2023-01-02",
				Created:     time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				Description: "Second test article.",
				Tags:        []string{"test", "second"},
			},
		}, nil
	case "feed_empty_desc":
		return []clippingsfeed.Metadata{
			{
				Title:       "No Description Article",
				Site:        "example.com",
				Source:      "https://example.com/nodesc",
				Author:      []string{"Test Author"},
				Published:   "2023-01-01",
				Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "",
				Tags:        []string{"test"},
			},
		}, nil
	default:
		return nil, nil
	}
}

// normalizeTimestamp replaces the timestamp in HTML with a fixed value for testing
func normalizeTimestamp(content string) string {
	// Replace any timestamp pattern with a fixed timestamp
	// This regex matches "Last updated: YYYY-MM-DD HH:MM:SS"
	return regexp.MustCompile(`Last updated: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`).
		ReplaceAllString(content, "Last updated: 2023-01-01 00:00:00")
}
