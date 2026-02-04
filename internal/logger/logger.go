package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Logger handles logging to file and console
type Logger struct {
	file    *os.File
	logger  *log.Logger
	verbose bool

	PagesCrawled int
	PagesFailed  int
	PagesSkipped int
}

// New creates a new Logger
func New(outputDir string, verbose bool) (*Logger, error) {
	logPath := filepath.Join(outputDir, "crawl.log")
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open log file
	file, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	l := &Logger{
		file:    file,
		logger:  log.New(file, "", log.LstdFlags),
		verbose: verbose,
	}

	return l, nil
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Println("INFO:", msg)
	if l.verbose {
		fmt.Println("INFO:", msg)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Println("WARNING:", msg)
	if l.verbose {
		fmt.Println("WARNING:", msg)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Println("ERROR:", msg)
	if l.verbose {
		fmt.Println("ERROR:", msg)
	}
}

// PageCrawled logs a successful page crawl
func (l *Logger) PageCrawled(url string) {
	l.PagesCrawled++
	l.Info("Crawled: %s", url)
}

// PageFailed logs a failed page crawl
func (l *Logger) PageFailed(url, reason string) {
	l.PagesFailed++
	l.Error("Failed to crawl %s: %s", url, reason)
}

// PageSkipped logs a skipped page
func (l *Logger) PageSkipped(url, reason string) {
	l.PagesSkipped++
	l.Info("Skipped %s: %s", url, reason)
}

// Summary returns crawl summary
func (l *Logger) Summary() string {
	total := l.PagesCrawled + l.PagesFailed + l.PagesSkipped
	return fmt.Sprintf(`
Crawl Summary:
  Total URLs processed: %d
  Successfully crawled: %d
  Failed: %d
  Skipped: %d
  Log file: %s`, total, l.PagesCrawled, l.PagesFailed, l.PagesSkipped, l.file.Name())
}

// Close closes the log file
func (l *Logger) Close() error {
	l.Info(l.Summary())
	return l.file.Close()
}
