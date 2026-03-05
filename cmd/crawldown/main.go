package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hex29a/crawldown/internal/config"
	"github.com/hex29a/crawldown/internal/crawler"
	"github.com/hex29a/crawldown/internal/localfile"
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
		Use:   "crawldown [URL or FILE]",
		Short: "CrawlDown - Website to Markdown crawler",
		Long: `CrawlDown crawls websites recursively and converts pages to clean, 
LLM-friendly markdown files organized in a structured folder hierarchy.

It can also process local HTML files — useful for sites with anti-bot 
protection (e.g., archive.ph). Save the page in your browser, then run 
CrawlDown on the saved file.`,
		Example: `  # Basic usage
  crawldown https://example.com

  # With custom options
  crawldown https://example.com --depth 5 --output ./docs

  # Exclude patterns
  crawldown https://example.com --exclude "*/tag/*" --exclude "*/category/*"

  # Allow additional domains
  crawldown https://example.com --allow-domain cdn.example.com

  # Process a local HTML file
  crawldown page.html --source-url https://original-site.com/page

  # Process all HTML files in a directory
  crawldown ./saved-pages/ --source-url https://original-site.com

  # Custom user agent
  crawldown https://example.com --user-agent "Mozilla/5.0 ..."`,
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
	rootCmd.Flags().StringVar(&cfg.UserAgent, "user-agent", config.DefaultUserAgent, "Custom User-Agent string for HTTP requests")
	rootCmd.Flags().StringVar(&cfg.SourceURL, "source-url", "", "Original source URL (used as metadata when processing local files)")
}

func run(cmd *cobra.Command, args []string) error {
	target := args[0]

	// Check if target is a local file/directory
	if localfile.IsLocalPath(target) {
		return runLocal(target)
	}

	cfg.StartURL = target
	return runCrawl()
}

func runLocal(target string) error {
	resolvedPath := localfile.ResolvePath(target)

	fmt.Printf("CrawlDown - Website to Markdown Crawler\n")
	fmt.Printf("Processing local: %s\n", resolvedPath)
	fmt.Printf("Output directory: %s\n", cfg.OutputDir)
	if cfg.SourceURL != "" {
		fmt.Printf("Source URL: %s\n", cfg.SourceURL)
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

	// Determine base URL for output organization
	baseURL := cfg.SourceURL
	if baseURL == "" {
		baseURL = "file://local"
	}

	// Create output manager
	out, err := output.New(cfg.OutputDir, baseURL)
	if err != nil {
		return fmt.Errorf("failed to create output manager: %w", err)
	}
	out.SetMergeMode(cfg.Merge, cfg.MergeBatchSize)

	// Create local file processor
	processor := localfile.New(log, out, cfg.SourceURL)

	// Check if target is a directory or file
	info, err := os.Stat(resolvedPath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", resolvedPath, err)
	}

	fmt.Println("Processing...")
	if info.IsDir() {
		if err := processor.ProcessDirectory(resolvedPath); err != nil {
			return fmt.Errorf("processing failed: %w", err)
		}
	} else {
		if err := processor.ProcessFile(resolvedPath); err != nil {
			return fmt.Errorf("processing failed: %w", err)
		}
	}

	// Flush any remaining buffered pages (merge mode)
	if err := out.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Processing completed!")
	summary := log.Summary()
	fmt.Println(summary)

	return nil
}

func runCrawl() error {
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
