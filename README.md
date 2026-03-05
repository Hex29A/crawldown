# CrawlDown

A powerful website-to-markdown crawler built in Go. CrawlDown recursively crawls websites and converts pages to clean, LLM-friendly markdown files with smart content extraction.

Perfect for creating documentation archives, feeding content to LLMs, or building searchable knowledge bases from websites.

## Features

- 🕷️ **Recursive Crawling** - Intelligently follows links with configurable depth limits
- 📝 **Clean Markdown Output** - Optimized for LLMs and human readability
- 🗂️ **Organized Structure** - Mirrors site hierarchy in your file system
- 🤖 **Respects Standards** - Honors robots.txt and implements rate limiting
- 🎯 **Smart Extraction** - Multiple content selectors for different site types
- 🌐 **Google Sites Support** - Special handling for Google Sites pages
- 📊 **Progress Tracking** - Detailed logging and crawl statistics
- 🔍 **URL Filtering** - Include/exclude patterns for precise control
- 🌍 **Multi-Domain** - Optionally crawl across multiple domains
- 📂 **Local HTML Files** - Process saved HTML files or directories (great for anti-bot sites)
- 🔧 **Custom User-Agent** - Configurable User-Agent for compatibility

## Installation

### Build from Source

```bash
git clone https://github.com/Hex29A/crawldown
cd crawldown
go build -o crawldown ./cmd/crawldown
```

### Using Go Install

```bash
go install github.com/Hex29A/crawldown/cmd/crawldown@latest
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/Hex29A/crawldown/releases).

## Quick Start

```bash
# Crawl a single page
./crawldown https://example.com --depth 0

# Crawl an entire site (depth 3)
./crawldown https://example.com

# Deep crawl with verbose output
./crawldown https://docs.example.com --depth 5 --verbose

# Crawl site with self-signed certificate
./crawldown https://example.com -k

# Process a local HTML file (e.g., saved from archive.ph)
./crawldown saved-page.html --source-url https://original-site.com/page

# Process all HTML files in a directory
./crawldown ./saved-pages/ --source-url https://original-site.com
```

## Usage

### Basic Command

```bash
crawldown <URL or FILE> [flags]
```

### Common Examples

**Crawl a blog with depth limit:**
```bash
./crawldown https://jvns.ca/ --depth 2 --output ./blog-archive
```

**Crawl documentation site:**
```bash
./crawldown https://go.dev/doc/ --depth 4 --delay 2s --verbose
```

**Crawl specific sections only:**
```bash
./crawldown https://example.com/docs \
  --include "*/docs/*" \
  --exclude "*/archive/*" \
  --depth 5
```

**Crawl with CDN support:**
```bash
./crawldown https://example.com \
  --allow-domain cdn.example.com \
  --allow-domain static.example.com \
  --depth 3
```

**Fast local testing:**
```bash
./crawldown http://localhost:8080 --delay 100ms --depth 2
```

**Crawl site with invalid/self-signed certificate:**
```bash
./crawldown https://internal.company.com -k --depth 3
```

## Command-Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--depth` | `-d` | Maximum crawl depth (0 = single page) | `3` |
| `--output` | `-o` | Output directory for markdown files | `./output` |
| `--delay` | | Delay between requests (e.g., 1s, 500ms) | `1s` |
| `--exclude` | | URL patterns to exclude (can specify multiple) | none |
| `--include` | | URL patterns to include (can specify multiple) | none |
| `--allow-domain` | | Additional domains to crawl | none |
| `--verbose` | `-v` | Enable detailed logging | `false` |
| `--insecure` | `-k` | Skip TLS certificate verification | `false` |
| `--user-agent` | | Custom User-Agent string for HTTP requests | browser-like |
| `--source-url` | | Original source URL (metadata for local file mode) | none |
### Depth Levels Explained

- `--depth 0` - Only crawl the specified URL
- `--depth 1` - Crawl specified URL + direct links
- `--depth 2` - Crawl up to 2 levels deep
- `--depth 3` - Default, good for most sites
- `--depth 5+` - Deep crawl for documentation sites

### URL Patterns

URL patterns support wildcards:
- `*/blog/*` - Match any URL with `/blog/` in the path
- `*/tag/*` - Exclude tag pages
- `*.pdf` - Match PDF files
- `/api/*` - Match API documentation

## Output Structure

CrawlDown creates an organized directory structure mirroring the website:

```
output/
├── crawl.log                    # Detailed crawl log
└── example.com/
    ├── index.md                 # Site map with all crawled URLs
    ├── home.md                  # Homepage content
    ├── about.md
    ├── blog/
    │   ├── post-1.md
    │   └── post-2.md
    └── docs/
        ├── getting-started.md
        └── api-reference.md
```

### Markdown Format

Each markdown file includes YAML frontmatter with metadata:

```markdown
---
source_url: https://example.com/docs/getting-started
title: Getting Started Guide
crawled_at: 2026-02-04T12:34:56Z
excerpt: "Quick introduction to getting started..."
---

# Getting Started Guide

Your clean markdown content here...
```

### Index File

The `index.md` file provides a complete site map:

```markdown
# Site Map: example.com

Crawled: 2026-02-04 12:34:56

## Pages (15)

- [Homepage](home.md) - https://example.com
- [About Us](about.md) - https://example.com/about
- [Getting Started](docs/getting-started.md) - https://example.com/docs/getting-started
...
```

## Content Extraction

CrawlDown uses intelligent content extraction with multiple strategies:

1. **Google Sites** - Special selector for `.tyJCtd` elements
2. **Semantic HTML** - Looks for `<main>`, `<article>` tags
3. **ARIA Roles** - Respects `role="main"` attributes
4. **Common Patterns** - Searches for `#content`, `.content` classes
5. **Fallback** - Uses `<body>` content if no main content found

This multi-layered approach ensures content is extracted cleanly from:
- Modern web apps
- Static site generators (Hugo, Jekyll)
- Documentation platforms
- Google Sites
- WordPress and other CMSs
- Custom HTML sites

## Real-World Examples

### Example 1: Personal Blog

```bash
./crawldown https://jvns.ca/ --depth 1
```

**Result:**
```
✓ Crawl completed!
  Total URLs: 1
  Successfully crawled: 1
  Output: output/jvns.ca/julia-evans.md (318 lines)
```

### Example 2: Google Sites

```bash
./crawldown https://www.hex29a.com --depth 2
```

**Result:**
```
✓ Crawl completed!
  Total URLs: 9
  Successfully crawled: 9
  Output: output/www.hex29a.com/ (9 pages)
```

### Example 3: Documentation Site

```bash
./crawldown https://developer.axis.com/vapix/ --depth 3 --delay 1s --verbose
```

**Result:**
```
✓ Crawl completed!
  Total URLs: 180
  Successfully crawled: 180
  Output: output/developer.axis.com/ (~103,000 lines of markdown)
```

## Performance & Best Practices

### Rate Limiting

Always use appropriate delays to avoid overwhelming servers:
- **Public sites**: `--delay 1s` (default)
- **Slower sites**: `--delay 2s` or more
- **Local development**: `--delay 100ms`

### Memory Usage

For large sites (100+ pages):
- Run with reasonable depth limits (3-5)
- Use `--exclude` to filter unnecessary pages
- Monitor with `--verbose` flag

### Common Exclusion Patterns

```bash
# Exclude common non-content pages
./crawldown https://example.com \
  --exclude "*/tag/*" \
  --exclude "*/category/*" \
  --exclude "*/author/*" \
  --exclude "*/page/*" \
  --exclude "*.pdf" \
  --exclude "*.jpg" \
  --exclude "*.png"
```

## Troubleshooting

### Anti-Bot Protection (429 Too Many Requests)

Some sites like archive.ph use aggressive anti-bot protection:
1. Save the page in your browser (Ctrl+S / Cmd+S)
2. Process the saved HTML file locally:
   ```bash
   ./crawldown saved-page.html --source-url https://archive.ph/kIMZN
   ```
3. Or try a custom User-Agent: `--user-agent "Mozilla/5.0 ..."`

### No Content Extracted

If pages have empty content:
1. Check if the site uses JavaScript rendering
2. Try `--verbose` to see extraction details
3. Verify the site's HTML structure
4. Some sites may block crawlers (403 Forbidden)

### Too Many Pages

If crawling too many pages:
1. Reduce `--depth`
2. Use `--include` patterns to target specific sections
3. Add `--exclude` patterns for unwanted pages

### Slow Crawling

If crawling is slow:
1. Reduce `--delay` (carefully)
2. Check your internet connection
3. Some sites may have rate limiting

## Development

### Project Structure

```
crawldown/
├── cmd/
│   └── crawldown/          # CLI entry point
│       └── main.go
├── internal/
│   ├── crawler/            # Web crawling logic
│   ├── extractor/          # Content extraction
│   ├── converter/          # HTML to Markdown
│   ├── localfile/          # Local HTML file processing
│   ├── output/             # File management
│   └── logger/             # Logging
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o crawldown ./cmd/crawldown
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see LICENSE file for details

## Credits

Built with:
- [Colly](https://github.com/gocolly/colly) - Web scraping framework
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - HTML conversion
- [Cobra](https://github.com/spf13/cobra) - CLI framework

---

**Note**: Always respect website terms of service and robots.txt. Use responsibly and ethically.
