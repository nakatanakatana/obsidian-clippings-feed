package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	clippingsfeed "github.com/nakatanakatana/obsidian-clippings-feed"
	"github.com/yuin/goldmark"
)

type FeedGenerator struct {
	config Config
	tmpDir string
	parser goldmark.Markdown
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

func (g *FeedGenerator) GenerateIndexHTML(filename string) error {
	metadata, err := g.scanMarkdownFiles()
	if err != nil {
		return fmt.Errorf("failed to scan markdown files: %w", err)
	}

	html := `<!DOCTYPE html>
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
        <h1>%s</h1>
        <p>%s</p>
    </div>
    
    <div class="feeds">
        <strong>Available feeds:</strong><br><br>
        <a href="/feed.rss">RSS</a>
        <a href="/feed.atom">Atom</a>
        <a href="/feed.json">JSON</a>
    </div>
    
    <div class="stats">
        <strong>Statistics:</strong> %d items found in directory: %s
    </div>
    
    <div class="content">
        <h2>Recent Items</h2>
        <ul class="items">%s</ul>
    </div>
    
    <div class="last-updated">
        Last updated: %s (refresh interval: %s)
    </div>
</body>
</html>`

	var itemsHTML string
	for _, meta := range metadata {
		tags := strings.Join(meta.Tags, ", ")
		if tags == "" {
			tags = "No tags"
		}

		itemHTML := fmt.Sprintf(`
            <li class="item">
                <div class="item-title"><a href="%s" target="_blank">%s</a></div>
                <div class="item-meta">Author: %s | Site: %s | Published: %s</div>
                <div class="item-desc">%s</div>
                <div class="item-tags">Tags: %s</div>
            </li>`,
			meta.Source, meta.Title, meta.Author, meta.Site, meta.Published, meta.Description, tags)
		itemsHTML += itemHTML
	}

	finalHTML := fmt.Sprintf(html,
		g.config.FeedTitle,
		g.config.FeedDesc,
		len(metadata),
		g.config.TargetDir,
		itemsHTML,
		time.Now().Format("2006-01-02 15:04:05"),
		g.config.RefreshInterval,
	)

	return os.WriteFile(filename, []byte(finalHTML), 0644)
}

func (g *FeedGenerator) StartPeriodicGeneration() {
	ticker := time.NewTicker(g.config.RefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := g.GenerateFeeds(); err != nil {
			log.Printf("Error during periodic feed generation: %v", err)
		}

		indexHTML := filepath.Join(g.tmpDir, "index.html")
		if err := g.GenerateIndexHTML(indexHTML); err != nil {
			log.Printf("Error during periodic index generation: %v", err)
		}
	}
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
