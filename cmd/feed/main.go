package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TargetDir       string        `env:"FEED_TARGET_DIR" envDefault:"./"`
	Port            string        `env:"FEED_PORT" envDefault:"8080"`
	FeedTitle       string        `env:"FEED_TITLE" envDefault:"Obsidian Clippings Feed"`
	FeedLink        string        `env:"FEED_LINK" envDefault:"http://localhost:8080"`
	FeedDesc        string        `env:"FEED_DESC" envDefault:"RSS feed from Obsidian clippings"`
	FeedAuthor      string        `env:"FEED_AUTHOR" envDefault:"Obsidian User"`
	MaxItems        int           `env:"FEED_MAX_ITEMS" envDefault:"50"`
	DebounceDelay   time.Duration `env:"FEED_DEBOUNCE_DELAY" envDefault:"10s"`
	HideDescription bool          `env:"FEED_HIDE_DESCRIPTION" envDefault:"false"`
}

func main() {
	// Configure structured logging
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	var config Config
	if err := env.Parse(&config); err != nil {
		slog.Error("Failed to parse config", "error", err)
		os.Exit(1)
	}

	tmpDir, err := os.MkdirTemp("", "obsidian-feed-*")
	if err != nil {
		slog.Error("Failed to create temp directory", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			slog.Warn("Failed to remove temp directory", "error", err, "tmpDir", tmpDir)
		}
	}()

	slog.Info("Created temp directory", "tmpDir", tmpDir)

	generator := NewFeedGenerator(config, tmpDir)

	if err := generator.GenerateFeeds(); err != nil {
		slog.Error("Failed to generate initial feeds", "error", err)
		os.Exit(1)
	}

	indexHTML := filepath.Join(tmpDir, "index.html")
	if err := generator.GenerateIndexHTML(indexHTML); err != nil {
		slog.Error("Failed to generate index.html", "error", err, "file", indexHTML)
		os.Exit(1)
	}

	if err := generator.StartFileWatcher(); err != nil {
		slog.Error("Failed to start file watcher", "error", err)
		os.Exit(1)
	}

	fs := http.FileServer(http.Dir(tmpDir))
	http.Handle("/", fs)

	slog.Info("Starting feed server",
		"port", config.Port,
		"watchDir", config.TargetDir,
		"serveDir", tmpDir,
		"debounceDelay", config.DebounceDelay,
		"hideDescription", config.HideDescription)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err, "port", config.Port)
		os.Exit(1)
	}
}
