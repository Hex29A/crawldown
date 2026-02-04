package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hex29a/crawldown/internal/config"
	"github.com/hex29a/crawldown/internal/converter"
	"github.com/hex29a/crawldown/internal/extractor"
	"github.com/hex29a/crawldown/internal/logger"
	"github.com/hex29a/crawldown/internal/output"
)

// Crawler handles website crawling
type Crawler struct {
	config  *config.Config
	logger  *logger.Logger
	output  *output.Manager
	visited map[string]bool
}

// New creates a new Crawler
func New(cfg *config.Config, log *logger.Logger, out *output.Manager) *Crawler {
	return &Crawler{
		config:  cfg,
		logger:  log,
		output:  out,
		visited: make(map[string]bool),
	}
}

// Start begins the crawling process
func (c *Crawler) Start() error {
	c.logger.Info("Starting crawl from %s", c.config.StartURL)
	c.logger.Info("Max depth: %d", c.config.MaxDepth)
	c.logger.Info("Allowed domains: %v", c.getAllowedDomains())

	collector := c.createCollector()

	// Handle HTML pages
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		c.processPage(e)
	})

	// Handle errors
	collector.OnError(func(r *colly.Response, err error) {
		c.logger.PageFailed(r.Request.URL.String(), err.Error())
	})

	// Start crawling
	if err := collector.Visit(c.config.StartURL); err != nil {
		return fmt.Errorf("failed to start crawl: %w", err)
	}

	// Wait for crawling to finish
	collector.Wait()

	// Create index
	c.logger.Info("Creating index file...")
	indexPath, err := c.output.CreateIndex()
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	c.logger.Info("Created index at %s", indexPath)

	return nil
}

func (c *Crawler) createCollector() *colly.Collector {
	options := []colly.CollectorOption{
		colly.MaxDepth(c.config.MaxDepth),
		colly.Async(false),
	}

	// Add allowed domains
	allowedDomains := c.getAllowedDomains()
	if len(allowedDomains) > 0 {
		options = append(options, colly.AllowedDomains(allowedDomains...))
	}

	collector := colly.NewCollector(options...)

	// Set rate limit
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       c.config.Delay,
		RandomDelay: c.config.Delay / 2,
	})

	// Set user agent
	collector.UserAgent = "CrawlDown/1.0 (Website to Markdown Crawler)"

	// Respect robots.txt
	collector.AllowURLRevisit = false

	// Handle links
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		// Check if should follow
		if c.shouldCrawl(absoluteURL) {
			e.Request.Visit(link)
		}
	})

	return collector
}

func (c *Crawler) processPage(e *colly.HTMLElement) {
	pageURL := e.Request.URL.String()

	// Check if already visited
	if c.visited[pageURL] {
		return
	}
	c.visited[pageURL] = true

	// Get HTML
	html, err := e.DOM.Html()
	if err != nil {
		c.logger.PageFailed(pageURL, fmt.Sprintf("failed to get HTML: %v", err))
		return
	}

	// Extract content
	title, contentHTML, excerpt, err := extractor.Extract(html, pageURL)
	if err != nil {
		c.logger.PageFailed(pageURL, fmt.Sprintf("extraction failed: %v", err))
		return
	}

	// Convert to markdown
	markdown, err := converter.ToMarkdown(contentHTML)
	if err != nil {
		c.logger.PageFailed(pageURL, fmt.Sprintf("conversion failed: %v", err))
		return
	}

	// Format with frontmatter
	finalMarkdown := converter.FormatWithFrontmatter(title, pageURL, markdown, excerpt)

	// Save to file
	filePath, err := c.output.SavePage(pageURL, title, finalMarkdown)
	if err != nil {
		c.logger.PageFailed(pageURL, fmt.Sprintf("save failed: %v", err))
		return
	}

	c.logger.PageCrawled(pageURL)
	c.logger.Info("Saved to %s", filePath)
}

func (c *Crawler) shouldCrawl(urlStr string) bool {
	// Parse URL
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if already visited
	normalizedURL := normalizeURL(urlStr)
	if c.visited[normalizedURL] {
		return false
	}

	// Check exclude patterns
	for _, pattern := range c.config.ExcludePatterns {
		if strings.Contains(urlStr, pattern) {
			c.logger.PageSkipped(urlStr, fmt.Sprintf("matches exclude pattern: %s", pattern))
			return false
		}
	}

	// Check include patterns (if specified)
	if len(c.config.IncludePatterns) > 0 {
		matches := false
		for _, pattern := range c.config.IncludePatterns {
			if strings.Contains(urlStr, pattern) {
				matches = true
				break
			}
		}
		if !matches {
			c.logger.PageSkipped(urlStr, "does not match include patterns")
			return false
		}
	}

	// Check if domain is allowed
	domain := parsed.Hostname()
	allowedDomains := c.getAllowedDomains()
	for _, allowed := range allowedDomains {
		if domain == allowed {
			return true
		}
	}

	return false
}

func (c *Crawler) getAllowedDomains() []string {
	domains := []string{}

	// Add starting domain
	if parsed, err := url.Parse(c.config.StartURL); err == nil {
		domains = append(domains, parsed.Hostname())
	}

	// Add configured domains
	domains = append(domains, c.config.AllowedDomains...)

	return domains
}

func normalizeURL(urlStr string) string {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove fragment
	parsed.Fragment = ""

	// Normalize path (remove trailing slash except for root)
	if parsed.Path != "/" && strings.HasSuffix(parsed.Path, "/") {
		parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	}

	return parsed.String()
}
