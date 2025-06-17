package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
		{
			name:     "filtered_items",
			testData: "filtered_metadata",
			config: Config{
				FeedTitle:     "Filtered Feed",
				FeedDesc:      "Feed with filtered items",
				FeedLink:      "http://example.com",
				FeedAuthor:    "Test Author",
				TargetDir:     "/test/dir",
				MaxItems:      50,
				DebounceDelay: 10 * time.Second,
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
	case "filtered_metadata":
		return []clippingsfeed.Metadata{
			{
				Title:       "Valid Article",
				Site:        "example.com",
				Source:      "https://example.com/valid",
				Author:      []string{"Valid Author"},
				Published:   "2023-01-01",
				Created:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "This article has both title and source",
				Tags:        []string{"valid"},
			},
			{
				Title:       "", // Empty title - should be filtered
				Site:        "example.com",
				Source:      "https://example.com/no-title",
				Author:      []string{"No Title Author"},
				Published:   "2023-01-02",
				Created:     time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				Description: "This article has no title",
				Tags:        []string{"invalid"},
			},
			{
				Title:       "No Source Article",
				Site:        "example.com",
				Source:      "", // Empty source - should be filtered
				Author:      []string{"No Source Author"},
				Published:   "2023-01-03",
				Created:     time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
				Description: "This article has no source",
				Tags:        []string{"invalid"},
			},
		}, nil
	default:
		return nil, nil
	}
}

func TestGenerateFeeds(t *testing.T) {
	tests := []struct {
		name        string
		testData    string
		config      Config
		expectError bool
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
		{
			name:     "limited_items",
			testData: "feed_multiple",
			config: Config{
				FeedTitle:     "Limited Feed",
				FeedDesc:      "Feed with item limit",
				FeedLink:      "http://example.com",
				FeedAuthor:    "Test Author",
				TargetDir:     "/test/dir",
				MaxItems:      1,
				DebounceDelay: 10 * time.Second,
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

			// Create a temporary markdown file for scanning
			testMarkdownDir := filepath.Join(tmpDir, "markdown")
			err := os.MkdirAll(testMarkdownDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test markdown directory: %v", err)
			}

			// Load test data and create markdown files
			metadata, err := loadTestData(t, tt.testData)
			if err != nil {
				t.Fatalf("Failed to load test data: %v", err)
			}

			// Create markdown files with YAML frontmatter
			for i, meta := range metadata {
				content := createMarkdownContent(meta)
				filename := filepath.Join(testMarkdownDir, fmt.Sprintf("test%d.md", i))
				err := os.WriteFile(filename, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test markdown file: %v", err)
				}
			}

			// Update generator's target directory to scan the test files
			generator.config.TargetDir = testMarkdownDir

			// Call GenerateFeeds
			err = generator.GenerateFeeds()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("GenerateFeeds failed: %v", err)
			}

			// Verify feed files were created
			feedFiles := []string{"feed.rss", "feed.atom", "feed.json"}
			for _, filename := range feedFiles {
				feedPath := filepath.Join(tmpDir, filename)
				if _, err := os.Stat(feedPath); os.IsNotExist(err) {
					t.Errorf("Feed file %s was not created", filename)
				}

				// Verify file has content
				content, err := os.ReadFile(feedPath)
				if err != nil {
					t.Errorf("Failed to read feed file %s: %v", filename, err)
				}
				if len(content) == 0 {
					t.Errorf("Feed file %s is empty", filename)
				}

				// Basic content validation
				contentStr := string(content)
				if !strings.Contains(contentStr, tt.config.FeedTitle) {
					t.Errorf("Feed file %s does not contain expected title", filename)
				}
			}
		})
	}
}

// createMarkdownContent creates markdown content with YAML frontmatter from metadata
func createMarkdownContent(meta clippingsfeed.Metadata) string {
	var content strings.Builder
	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("title: \"%s\"\n", meta.Title))
	content.WriteString(fmt.Sprintf("site: \"%s\"\n", meta.Site))
	content.WriteString(fmt.Sprintf("source: \"%s\"\n", meta.Source))

	if len(meta.Author) > 0 {
		content.WriteString("author:\n")
		for _, author := range meta.Author {
			content.WriteString(fmt.Sprintf("  - \"%s\"\n", author))
		}
	}

	content.WriteString(fmt.Sprintf("published: \"%s\"\n", meta.Published))
	content.WriteString(fmt.Sprintf("created: %s\n", meta.Created.Format(time.RFC3339)))
	content.WriteString(fmt.Sprintf("description: \"%s\"\n", meta.Description))

	if len(meta.Tags) > 0 {
		content.WriteString("tags:\n")
		for _, tag := range meta.Tags {
			content.WriteString(fmt.Sprintf("  - \"%s\"\n", tag))
		}
	}

	content.WriteString("---\n\n")
	content.WriteString("# " + meta.Title + "\n\n")
	content.WriteString(meta.Description + "\n")

	return content.String()
}

// normalizeTimestamp replaces the timestamp in HTML with a fixed value for testing
func normalizeTimestamp(content string) string {
	// Replace any timestamp pattern with a fixed timestamp
	// This regex matches "Last updated: YYYY-MM-DD HH:MM:SS"
	return regexp.MustCompile(`Last updated: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`).
		ReplaceAllString(content, "Last updated: 2023-01-01 00:00:00")
}
