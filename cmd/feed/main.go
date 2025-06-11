package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TargetDir     string        `env:"FEED_TARGET_DIR" envDefault:"./"`
	Port          string        `env:"FEED_PORT" envDefault:"8080"`
	FeedTitle     string        `env:"FEED_TITLE" envDefault:"Obsidian Clippings Feed"`
	FeedLink      string        `env:"FEED_LINK" envDefault:"http://localhost:8080"`
	FeedDesc      string        `env:"FEED_DESC" envDefault:"RSS feed from Obsidian clippings"`
	FeedAuthor    string        `env:"FEED_AUTHOR" envDefault:"Obsidian User"`
	DebounceDelay time.Duration `env:"FEED_DEBOUNCE_DELAY" envDefault:"10s"`
}

func main() {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "obsidian-feed-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("Failed to remove temp directory: %v", err)
		}
	}()

	log.Printf("Created temp directory: %s", tmpDir)

	generator := NewFeedGenerator(config, tmpDir)

	if err := generator.GenerateFeeds(); err != nil {
		log.Fatalf("Failed to generate initial feeds: %v", err)
	}

	indexHTML := filepath.Join(tmpDir, "index.html")
	if err := generator.GenerateIndexHTML(indexHTML); err != nil {
		log.Fatalf("Failed to generate index.html: %v", err)
	}

	if err := generator.StartFileWatcher(); err != nil {
		log.Fatalf("Failed to start file watcher: %v", err)
	}

	fs := http.FileServer(http.Dir(tmpDir))
	http.Handle("/", fs)

	log.Printf("Starting feed server on port %s", config.Port)
	log.Printf("Watching directory: %s", config.TargetDir)
	log.Printf("Serving files from: %s", tmpDir)
	log.Printf("Debounce delay: %s", config.DebounceDelay)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
