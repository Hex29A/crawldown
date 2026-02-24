package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hex29a/crawldown/internal/config"
	"github.com/hex29a/crawldown/internal/crawler"
	"github.com/hex29a/crawldown/internal/logger"
	"github.com/hex29a/crawldown/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfg     *config.Config
	rootCmd *cobra.Command
)

func init() {
	cfg = config.New()

	rootCmd = &cobra.Command{
		Use:   "crawldown [URL]",
		Short: "CrawlDown - Website to Markdown crawler",
		Long: `CrawlDown crawls websites recursively and converts pages to clean, 
LLM-friendly markdown files organized in a structured folder hierarchy.`,
		Example: `  # Basic usage
  crawldown https://example.com

  # With custom options
  crawldown https://example.com --depth 5 --output ./docs

  # Exclude patterns
  crawldown https://example.com --exclude "*/tag/*" --exclude "*/category/*"

  # Allow additional domains
  crawldown https://example.com --allow-domain cdn.example.com`,
		Args: cobra.ExactArgs(1),
		RunE: run,
	}

	// Flags
	rootCmd.Flags().IntVarP(&cfg.MaxDepth, "depth", "d", 3, "Maximum crawl depth")
	rootCmd.Flags().StringVarP(&cfg.OutputDir, "output", "o", "output", "Output directory")
	rootCmd.Flags().DurationVar(&cfg.Delay, "delay", time.Second, "Delay between requests")
	rootCmd.Flags().StringSliceVar(&cfg.ExcludePatterns, "exclude", []string{}, "URL patterns to exclude")
	rootCmd.Flags().StringSliceVar(&cfg.IncludePatterns, "include", []string{}, "URL patterns to include")
	rootCmd.Flags().StringSliceVar(&cfg.AllowedDomains, "allow-domain", []string{}, "Additional domains to crawl")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.Flags().BoolVarP(&cfg.InsecureTLS, "insecure", "k", false, "Skip TLS certificate verification")
	rootCmd.Flags().BoolVarP(&cfg.Merge, "merge", "m", false, "Merge all pages into bundled markdown files (optimized for NotebookLM)")
	rootCmd.Flags().IntVar(&cfg.MergeBatchSize, "merge-batch-size", 50, "Number of pages per merged file")
}

func run(cmd *cobra.Command, args []string) error {
	cfg.StartURL = args[0]

	// Print startup info
	fmt.Printf("CrawlDown - Website to Markdown Crawler\n")
	fmt.Printf("Starting crawl of: %s\n", cfg.StartURL)
	fmt.Printf("Max depth: %d\n", cfg.MaxDepth)
	fmt.Printf("Output directory: %s\n", cfg.OutputDir)
	fmt.Printf("Request delay: %s\n", cfg.Delay)
	if len(cfg.AllowedDomains) > 0 {
		fmt.Printf("Additional domains: %v\n", cfg.AllowedDomains)
	}
	if len(cfg.ExcludePatterns) > 0 {
		fmt.Printf("Exclude patterns: %v\n", cfg.ExcludePatterns)
	}
	if len(cfg.IncludePatterns) > 0 {
		fmt.Printf("Include patterns: %v\n", cfg.IncludePatterns)
	}
	if cfg.Merge {
		fmt.Printf("Merge mode: enabled (batch size: %d)\n", cfg.MergeBatchSize)
	}
	fmt.Println()

	// Create logger
	log, err := logger.New(cfg.OutputDir, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer log.Close()

	// Create output manager
	out, err := output.New(cfg.OutputDir, cfg.StartURL)
	if err != nil {
		return fmt.Errorf("failed to create output manager: %w", err)
	}
	out.SetMergeMode(cfg.Merge, cfg.MergeBatchSize)

	// Create crawler
	c := crawler.New(cfg, log, out)

	// Start crawling
	fmt.Println("Crawling...")
	if err := c.Start(); err != nil {
		return fmt.Errorf("crawl failed: %w", err)
	}

	// Flush any remaining buffered pages (merge mode)
	if err := out.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	// Print summary
	fmt.Println()
	fmt.Println("✓ Crawl completed!")
	summary := log.Summary()
	fmt.Println(summary)

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
