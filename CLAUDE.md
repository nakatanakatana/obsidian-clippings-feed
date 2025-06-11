# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a dual-purpose project for Obsidian workflow automation:

1. **Primary**: Go-based HTTP feed server with real-time file watching for Obsidian clippings
2. **Secondary**: TypeScript/Deno-based LiveSync bridge service for peer-to-peer Obsidian synchronization

The main Go module (`github.com/nakatanakatana/obsidian-clippings-feed`) provides a complete feed server that monitors markdown files, extracts metadata from YAML frontmatter, and serves RSS/Atom/JSON feeds via HTTP with automatic regeneration on file changes.

## Development Commands

### Build and Test
```bash
# Build all CLI tools
go build -o dist/ ./cmd/...

# Run feed server (development)
./dist/feed

# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -v -coverprofile=coverage.txt -covermode=atomic

# Run specific test
go test . -run TestParseMeta -v

# Update golden files
go test ./... -update

# Format code
go fmt ./...

# Lint code
go tool golangci-lint run
```

## Post-Change Verification

**IMPORTANT**: After making any source code changes, always run these commands in sequence to ensure code quality:

```bash
# 1. Format code
go fmt ./...

# 2. Run tests
go test ./... -v

# 3. Run linter
go tool golangci-lint run

# 4. Build to verify compilation
go build -o dist/ ./cmd/...
```

If any step fails, fix the issues before committing. This workflow ensures:
- Consistent code formatting
- All tests pass
- No linting violations  
- Code compiles successfully

### Docker Operations
```bash
# Build main feed server image
docker build -t obsidian-feed .

# Run feed server in container
docker run -p 8080:8080 -v $(pwd):/data -e FEED_TARGET_DIR=/data obsidian-feed

# LiveSync bridge (separate service)
cd livesync-bridge
docker-compose up
```

## Core Architecture

### Feed Server with File Watching
The project provides a complete HTTP feed server with real-time monitoring:
1. **Web Clipper** (`clipper-template.json`) → captures web content with structured metadata
2. **File Watcher** (`fsnotify`) → monitors markdown files for changes recursively
3. **Parser** (`parser.go`) → extracts `Metadata` struct from markdown YAML frontmatter  
4. **Feed Generator** (`feed.go`) → converts `[]Metadata` to RSS/Atom/JSON feeds using gorilla/feeds
5. **HTTP Server** (`cmd/feed/`) → serves feeds and web interface with automatic regeneration

### Key Data Structure
```go
type Metadata struct {
    Title       string    `json:"title"`
    Site        string    `json:"site"`
    Source      string    `json:"source"`
    Author      []string  `json:"author"`    // Supports multiple authors
    Published   string    `json:"published"`
    Created     time.Time `json:"created"`
    Description string    `json:"description"`
    Tags        []string  `json:"tags"`
}
```

### Feed Server Configuration
Environment variables for feed server configuration:
```bash
FEED_TARGET_DIR=./              # Directory to watch for markdown files
FEED_PORT=8080                  # HTTP server port
FEED_TITLE="Obsidian Clippings Feed"  # Feed title
FEED_LINK="http://localhost:8080"     # Feed base URL
FEED_DESC="RSS feed from Obsidian clippings"  # Feed description
FEED_AUTHOR="Obsidian User"     # Default feed author
FEED_MAX_ITEMS=50               # Maximum items per feed
FEED_DEBOUNCE_DELAY=10s         # File change debounce delay
```

### Server Endpoints
- `/` - Web dashboard with feed links and statistics
- `/feed.rss` - RSS 2.0 format feed
- `/feed.atom` - Atom 1.0 format feed  
- `/feed.json` - JSON Feed format

### Testing Strategy
- **Golden file testing** with `gotest.tools/v3/golden` for deterministic output verification
- Test data in `testdata/` directory with Atom feed format as standard
- Comprehensive test coverage:
  - Single/multiple item feeds (`feed_single`, `feed_multiple`)
  - Empty description handling (`feed_empty_desc`)
  - Item limiting functionality (`feed_limited`)
  - Metadata parsing validation (`success`)
- URL/domain testing uses `example.com` for consistency
- Feed generation tests validate end-to-end functionality from metadata to Atom XML

### Dependencies
- **goldmark + goldmark-meta**: Markdown parsing with YAML frontmatter support
- **gorilla/feeds**: RSS/Atom/JSON feed generation (supports RSS, Atom, JSON formats)
- **fsnotify/fsnotify**: File system monitoring for real-time updates
- **caarlos0/env**: Environment variable parsing with defaults
- **gotest.tools/v3**: Enhanced testing utilities with golden file support

## LiveSync Bridge Component
Located in `livesync-bridge/`, this is a separate Deno/TypeScript application for Obsidian peer-to-peer synchronization. It uses a hub-based architecture with CouchDB integration and includes a Svelte-based web interface.

## Configuration Files
- `clipper-template.json`: Web clipper browser extension template for consistent metadata capture
- `livesync-bridge/dat/config.sample.json`: Configuration template for LiveSync bridge service
- `.github/workflows/`: CI/CD pipelines for Go linting, testing, and Docker publishing
