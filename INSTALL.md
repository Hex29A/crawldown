# CrawlDown-Go Installation Guide

## Prerequisites

### Install Go

#### Ubuntu/Debian
```bash
sudo snap install go --classic
# or
sudo apt install golang-go
```

#### macOS
```bash
brew install go
```

#### Windows
Download from https://go.dev/dl/

## Build Instructions

```bash
# Navigate to project
cd /home/martincr/Dev/Personal/CrawlDown-Go

# Download dependencies
go mod tidy

# Build the binary
go build -o crawldown ./cmd/crawldown

# Or install globally
go install ./cmd/crawldown
```

## Quick Start

```bash
# After building, run:
./crawldown https://example.com

# With options:
./crawldown https://example.com \
  --depth 5 \
  --output ./docs \
  --verbose
```

## Development

```bash
# Run without building
go run ./cmd/crawldown https://example.com

# Run tests
go test ./...

# Format code
go fmt ./...
```

## Cross-Compilation

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o crawldown-linux ./cmd/crawldown

# For macOS
GOOS=darwin GOARCH=amd64 go build -o crawldown-macos ./cmd/crawldown

# For Windows
GOOS=windows GOARCH=amd64 go build -o crawldown.exe ./cmd/crawldown
```
