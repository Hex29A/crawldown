# CrawlDown - Production Status

## ✅ Ready for Production

**Version**: 1.0.0  
**Language**: Go 1.21+  
**Binary Size**: 18MB  
**Status**: Fully Tested & Working

## Features

- [x] Recursive website crawling with depth limits
- [x] robots.txt respect and rate limiting
- [x] Smart content extraction (readability algorithm)
- [x] Clean markdown conversion
- [x] LLM-optimized output with frontmatter
- [x] Organized folder structure
- [x] Index file generation
- [x] URL normalization and deduplication
- [x] Domain filtering with allow/exclude patterns
- [x] Progress logging to file
- [x] Comprehensive error handling
- [x] Version metadata via --version / version command

## Tested On

- ✅ example.com (simple site)
- ✅ developer.axis.com/vapix/ (complex documentation, 180 pages)

## Quick Start

\`\`\`bash
# Build
go build -o crawldown ./cmd/crawldown

# Basic usage
./crawldown https://example.com

# Full featured
./crawldown https://example.com \\
  --depth 3 \\
  --output ./docs \\
  --delay 1s \\
  --verbose
\`\`\`

## Project History

- **2026-02-04**: Initial Python implementation
- **2026-02-04**: Go reimplementation completed
- **2026-02-04**: Python version archived, Go version promoted to primary

## Architecture

\`\`\`
cmd/crawldown/     - CLI entry point
internal/
  ├── config/      - Configuration management
  ├── crawler/     - Web crawling engine
  ├── extractor/   - HTML content extraction
  ├── converter/   - Markdown conversion
  ├── output/      - File management
  └── logger/      - Logging system
\`\`\`

## Dependencies

- gocolly/colly - Web crawling + robots.txt
- PuerkitoBio/goquery - HTML parsing
- JohannesKaufmann/html-to-markdown - Markdown conversion
- spf13/cobra - CLI framework

All dependencies vendored in go.mod.

## Development

\`\`\`bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Update dependencies
go mod tidy

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o crawldown-linux ./cmd/crawldown
GOOS=darwin GOARCH=amd64 go build -o crawldown-macos ./cmd/crawldown
GOOS=windows GOARCH=amd64 go build -o crawldown.exe ./cmd/crawldown
\`\`\`
