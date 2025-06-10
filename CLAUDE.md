# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a dual-purpose project for Obsidian workflow automation:

1. **Primary**: Go-based metadata parser and RSS feed generator for Obsidian clippings
2. **Secondary**: TypeScript/Deno-based LiveSync bridge service for peer-to-peer Obsidian synchronization

The main Go module (`github.com/nakatanakatana/obsidian-clippings-feed`) processes markdown files with YAML frontmatter to extract structured metadata and generate RSS/Atom feeds.

## Development Commands

### Build and Test
```bash
# Build the CLI tool
go build -o dist/ ./cmd/...

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
# Build main image
docker build -t obsidian-feed .

# LiveSync bridge
cd livesync-bridge
docker-compose up
```

## Core Architecture

### Metadata Processing Pipeline
The project follows this data flow:
1. **Web Clipper** (`clipper-template.json`) → captures web content with structured metadata
2. **Parser** (`parser.go`) → extracts `Metadata` struct from markdown YAML frontmatter  
3. **Feed Generator** (`feed.go`) → converts `[]Metadata` to RSS/Atom feeds using gorilla/feeds
4. **CLI Tool** (`cmd/obsidian-feed/main.go`) → orchestrates the pipeline (currently placeholder)

### Key Data Structure
```go
type Metadata struct {
    Title       string    `json:"title"`
    Site        string    `json:"site"`
    Source      string    `json:"source"`
    Author      string    `json:"author"`
    Published   string    `json:"published"`
    Created     time.Time `json:"created"`
    Description string    `json:"description"`
    Tags        []string  `json:"tags"`
}
```

### Testing Strategy
- **Golden file testing** with `gotest.tools/v3/golden` for deterministic output verification
- Test data in `testdata/` directory with Atom feed format as standard
- URL/domain testing uses `example.com` for consistency
- Feed generation tests validate end-to-end functionality from metadata to Atom XML

### Dependencies
- **goldmark + goldmark-meta**: Markdown parsing with YAML frontmatter support
- **gorilla/feeds**: RSS/Atom feed generation (supports RSS, Atom, JSON formats)
- **gotest.tools/v3**: Enhanced testing utilities with golden file support

## LiveSync Bridge Component
Located in `livesync-bridge/`, this is a separate Deno/TypeScript application for Obsidian peer-to-peer synchronization. It uses a hub-based architecture with CouchDB integration and includes a Svelte-based web interface.

## Configuration Files
- `clipper-template.json`: Web clipper browser extension template for consistent metadata capture
- `livesync-bridge/dat/config.sample.json`: Configuration template for LiveSync bridge service
- `.github/workflows/`: CI/CD pipelines for Go linting, testing, and Docker publishing
