package localfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hex29a/crawldown/internal/converter"
	"github.com/hex29a/crawldown/internal/extractor"
	"github.com/hex29a/crawldown/internal/logger"
	"github.com/hex29a/crawldown/internal/output"
)

// Processor handles local HTML file processing
type Processor struct {
	logger    *logger.Logger
	output    *output.Manager
	sourceURL string
}

// New creates a new local file Processor
func New(log *logger.Logger, out *output.Manager, sourceURL string) *Processor {
	return &Processor{
		logger:    log,
		output:    out,
		sourceURL: sourceURL,
	}
}

// ProcessFile reads a local HTML file and converts it to markdown
func (p *Processor) ProcessFile(filePath string) error {
	p.logger.Info("Processing local file: %s", filePath)

	// Read the HTML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	html := string(data)

	// Use the source URL if provided, otherwise use the file path
	pageURL := p.sourceURL
	if pageURL == "" {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			absPath = filePath
		}
		pageURL = "file://" + absPath
	}

	// Extract content
	title, contentHTML, excerpt, err := extractor.Extract(html, pageURL)
	if err != nil {
		return fmt.Errorf("extraction failed for %s: %w", filePath, err)
	}

	// Convert to markdown
	markdown, err := converter.ToMarkdown(contentHTML)
	if err != nil {
		return fmt.Errorf("conversion failed for %s: %w", filePath, err)
	}

	if p.output.IsMergeMode() {
		if err := p.output.BufferPage(pageURL, title, markdown); err != nil {
			return fmt.Errorf("buffer failed for %s: %w", filePath, err)
		}
		p.logger.Info("Buffered for merge: %s", title)
	} else {
		finalMarkdown := converter.FormatWithFrontmatter(title, pageURL, markdown, excerpt)

		savePath, err := p.output.SavePage(pageURL, title, finalMarkdown)
		if err != nil {
			return fmt.Errorf("save failed for %s: %w", filePath, err)
		}
		p.logger.PageCrawled(pageURL)
		p.logger.Info("Saved to %s", savePath)
	}

	return nil
}

// ProcessDirectory reads all HTML files in a directory and converts them
func (p *Processor) ProcessDirectory(dirPath string) error {
	p.logger.Info("Processing local directory: %s", dirPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".html" && ext != ".htm" {
			continue
		}

		fullPath := filepath.Join(dirPath, name)
		if err := p.ProcessFile(fullPath); err != nil {
			p.logger.PageFailed(fullPath, err.Error())
			continue
		}
		count++
	}

	if count == 0 {
		return fmt.Errorf("no HTML files found in %s", dirPath)
	}

	p.logger.Info("Processed %d HTML files from %s", count, dirPath)
	return nil
}

// IsLocalPath returns true if the given string looks like a local file or directory path
func IsLocalPath(path string) bool {
	// Check for file:// scheme
	if strings.HasPrefix(path, "file://") {
		return true
	}

	// Check if it's an existing file or directory
	if info, err := os.Stat(path); err == nil {
		_ = info
		return true
	}

	// Check common path patterns (starts with /, ./, ../, ~/)
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") || strings.HasPrefix(path, "~/") {
		return true
	}

	// Check if it has an HTML extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".html" || ext == ".htm" {
		return true
	}

	return false
}

// ResolvePath resolves a path, handling file:// prefix
func ResolvePath(path string) string {
	if strings.HasPrefix(path, "file://") {
		return strings.TrimPrefix(path, "file://")
	}
	return path
}
