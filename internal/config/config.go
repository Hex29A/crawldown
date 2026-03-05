package config

import "time"

// Config holds crawl configuration
type Config struct {
	StartURL        string
	MaxDepth        int
	OutputDir       string
	Delay           time.Duration
	AllowedDomains  []string
	ExcludePatterns []string
	IncludePatterns []string
	Verbose         bool
	InsecureTLS     bool
	Merge           bool
	MergeBatchSize  int
	UserAgent       string
	SourceURL       string // Original URL for local file mode metadata
}

// DefaultUserAgent is the default User-Agent string used for requests
const DefaultUserAgent = "Mozilla/5.0 (compatible; CrawlDown/1.0; +https://github.com/hex29a/crawldown)"

// New creates a new Config with defaults
func New() *Config {
	return &Config{
		MaxDepth:       3,
		OutputDir:      "output",
		Delay:          time.Second,
		Verbose:        false,
		Merge:          false,
		MergeBatchSize: 50,
		UserAgent:      DefaultUserAgent,
	}
}
