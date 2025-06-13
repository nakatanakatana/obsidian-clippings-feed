package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	clippingsfeed "github.com/nakatanakatana/obsidian-clippings-feed"
	"github.com/yuin/goldmark"
)

type FeedGenerator struct {
	config        Config
	tmpDir        string
	parser        goldmark.Markdown
	watcher       *fsnotify.Watcher
	debounceTimer *time.Timer
}

func NewFeedGenerator(config Config, tmpDir string) *FeedGenerator {
	return &FeedGenerator{
		config: config,
		tmpDir: tmpDir,
		parser: clippingsfeed.CreateParser(),
	}
}

func (g *FeedGenerator) GenerateFeeds() error {
	metadata, err := g.scanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("failed to scan markdown files: %w", err)
	}

	feedConfig := clippingsfeed.FeedConfig{
		Title:       g.config.FeedTitle,
		Link:        g.config.FeedLink,
		Description: g.config.FeedDesc,
		Author:      g.config.FeedAuthor,
		Created:     time.Now(),
		MaxItems:    g.config.MaxItems,
	}

	feed, err := clippingsfeed.GenerateFeed(metadata, feedConfig)
	if err != nil {
		return fmt.Errorf("failed to generate feed: %w", err)
	}

	if err := clippingsfeed.WriteFeedToFile(feed, filepath.Join(g.tmpDir, "feed.rss"), "rss"); err != nil {
		return fmt.Errorf("failed to write RSS feed: %w", err)
	}

	if err := clippingsfeed.WriteFeedToFile(feed, filepath.Join(g.tmpDir, "feed.atom"), "atom"); err != nil {
		return fmt.Errorf("failed to write Atom feed: %w", err)
	}

	if err := clippingsfeed.WriteFeedToFile(feed, filepath.Join(g.tmpDir, "feed.json"), "json"); err != nil {
		return fmt.Errorf("failed to write JSON feed: %w", err)
	}

	log.Printf("Generated feeds with %d items", len(metadata))
	return nil
}

// Template data structure for HTML rendering
type IndexTemplateData struct {
	FeedTitle   string
	FeedDesc    string
	ItemCount   int
	TargetDir   string
	Items       []IndexItem
	LastUpdated string
	UpdateMode  string
}

type IndexItem struct {
	Title       string
	Source      string
	Authors     string
	Site        string
	Published   string
	Description string
	Tags        string
}

const indexHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Obsidian Clippings Feed</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { margin-bottom: 30px; }
        .feeds { margin-bottom: 30px; }
        .feeds a { margin-right: 15px; padding: 5px 10px; background: #007cba; color: white; text-decoration: none; border-radius: 3px; }
        .stats { margin-bottom: 30px; color: #666; }
        .items { list-style: none; padding: 0; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .item-title { font-weight: bold; margin-bottom: 5px; }
        .item-meta { color: #666; font-size: 0.9em; margin-bottom: 10px; }
        .item-desc { margin-bottom: 10px; }
        .item-tags { font-size: 0.8em; color: #999; }
        .last-updated { color: #999; font-size: 0.9em; margin-top: 30px; text-align: center; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.FeedTitle}}</h1>
        <p>{{.FeedDesc}}</p>
    </div>
    
    <div class="feeds">
        <strong>Available feeds:</strong><br><br>
        <a href="/feed.rss">RSS</a>
        <a href="/feed.atom">Atom</a>
        <a href="/feed.json">JSON</a>
    </div>
    
    <div class="stats">
        <strong>Statistics:</strong> {{.ItemCount}} items found in directory: {{.TargetDir}}
    </div>
    
    <div class="content">
        <h2>Recent Items</h2>
        <ul class="items">
        {{range .Items}}
            <li class="item">
                <div class="item-title"><a href="{{.Source}}" target="_blank">{{.Title}}</a></div>
                <div class="item-meta">Author(s): {{.Authors}} | Site: {{.Site}} | Published: {{.Published}}</div>
                <div class="item-desc">{{.Description}}</div>
                <div class="item-tags">Tags: {{.Tags}}</div>
            </li>
        {{end}}
        </ul>
    </div>
    
    <div class="last-updated">
        Last updated: {{.LastUpdated}} (update mode: {{.UpdateMode}})
    </div>
</body>
</html>`

func (g *FeedGenerator) GenerateIndexHTML(filename string) error {
	metadata, err := g.scanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("failed to scan markdown files: %w", err)
	}

	return g.generateIndexHTMLFromMetadata(filename, metadata)
}

func (g *FeedGenerator) generateIndexHTMLFromMetadata(filename string, metadata []clippingsfeed.Metadata) error {
	// Create template
	tmpl, err := template.New("index").Parse(indexHTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	items := make([]IndexItem, len(metadata))
	for i, meta := range metadata {
		tags := strings.Join(meta.Tags, ", ")
		if tags == "" {
			tags = "No tags"
		}

		authors := strings.Join(meta.Author, ", ")
		if authors == "" {
			authors = "Unknown"
		}

		items[i] = IndexItem{
			Title:       meta.Title,
			Source:      meta.Source,
			Authors:     authors,
			Site:        meta.Site,
			Published:   meta.Published,
			Description: meta.Description,
			Tags:        tags,
		}
	}

	data := IndexTemplateData{
		FeedTitle:   g.config.FeedTitle,
		FeedDesc:    g.config.FeedDesc,
		ItemCount:   len(metadata),
		TargetDir:   g.config.TargetDir,
		Items:       items,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		UpdateMode:  "file watcher",
	}

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() //nolint:errcheck

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func (g *FeedGenerator) StartFileWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	g.watcher = watcher

	if err := g.addWatchesRecursively(g.config.TargetDir); err != nil {
		return fmt.Errorf("failed to add watches: %w", err)
	}

	go g.watchLoop()
	log.Printf("File watcher started for directory: %s", g.config.TargetDir)
	return nil
}

func (g *FeedGenerator) addWatchesRecursively(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if err := g.watcher.Add(path); err != nil {
				log.Printf("Warning: failed to watch directory %s: %v", path, err)
				return nil
			}
			log.Printf("Watching directory: %s", path)
		}
		return nil
	})
}

func (g *FeedGenerator) watchLoop() {
	defer func() {
		if err := g.watcher.Close(); err != nil {
			log.Printf("Error closing file watcher: %v", err)
		}
	}()

	for {
		select {
		case event, ok := <-g.watcher.Events:
			if !ok {
				return
			}

			if g.shouldProcessEvent(event) {
				g.debouncedRegenerate()
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if stat, err := os.Stat(event.Name); err == nil && stat.IsDir() {
					if err := g.watcher.Add(event.Name); err != nil {
						log.Printf("Warning: failed to watch new directory %s: %v", event.Name, err)
					} else {
						log.Printf("Added watch for new directory: %s", event.Name)
					}
				}
			}

		case err, ok := <-g.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)
		}
	}
}

func (g *FeedGenerator) shouldProcessEvent(event fsnotify.Event) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
		return false
	}

	if strings.HasSuffix(strings.ToLower(event.Name), ".md") {
		log.Printf("Detected change in markdown file: %s", event.Name)
		return true
	}

	return false
}

func (g *FeedGenerator) debouncedRegenerate() {
	if g.debounceTimer != nil {
		g.debounceTimer.Stop()
	}

	g.debounceTimer = time.AfterFunc(g.config.DebounceDelay, func() {
		log.Printf("Regenerating feeds due to file changes...")

		if err := g.GenerateFeeds(); err != nil {
			log.Printf("Error during feed regeneration: %v", err)
		}

		indexHTML := filepath.Join(g.tmpDir, "index.html")
		if err := g.GenerateIndexHTML(indexHTML); err != nil {
			log.Printf("Error during index regeneration: %v", err)
		}

		log.Printf("Feed regeneration completed")
	})
}

func (g *FeedGenerator) scanMarkdownFiles() ([]clippingsfeed.Metadata, error) {
	var metadata []clippingsfeed.Metadata

	err := filepath.WalkDir(g.config.TargetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		meta, err := clippingsfeed.ParseMeta(g.parser, string(content))
		if err != nil {
			log.Printf("Error parsing metadata from %s: %v", path, err)
			return nil
		}

		if meta.Title == "" {
			meta.Title = filepath.Base(path)
		}
		if meta.Created.IsZero() {
			if stat, err := os.Stat(path); err == nil {
				meta.Created = stat.ModTime()
			}
		}

		metadata = append(metadata, *meta)
		return nil
	})

	return metadata, err
}
