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
}

// New creates a new Config with defaults
func New() *Config {
	return &Config{
		MaxDepth:       3,
		OutputDir:      "output",
		Delay:          time.Second,
		Verbose:        false,
		Merge:          false,
		MergeBatchSize: 50,
	}
}
