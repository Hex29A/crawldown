package output

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Manager manages output files and organization
type Manager struct {
	outputDir      string
	siteDir        string
	baseURL        string
	domain         string
	filenameCounts map[string]int
	pages          []PageInfo

	// Merge mode fields
	mergeMode       bool
	mergeBatchSize  int
	mergeBuffer     []MergedPage
	mergePartNum    int
	mergeTotalPages int
}

// PageInfo stores page metadata
type PageInfo struct {
	URL   string
	Title string
	File  string
}

// MergedPage holds buffered page content for merge mode
type MergedPage struct {
	URL      string
	Title    string
	Markdown string
}

// New creates a new output Manager
func New(outputDir, baseURL string) (*Manager, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	domain := parsed.Hostname()
	siteDir := filepath.Join(outputDir, domain)

	if err := os.MkdirAll(siteDir, 0755); err != nil {
		return nil, err
	}

	return &Manager{
		outputDir:      outputDir,
		siteDir:        siteDir,
		baseURL:        baseURL,
		domain:         domain,
		filenameCounts: make(map[string]int),
		pages:          []PageInfo{},
	}, nil
}

// SavePage saves markdown content to file
func (m *Manager) SavePage(pageURL, title, markdown string) (string, error) {
	filename := m.generateFilename(pageURL, title)
	filePath := filepath.Join(m.siteDir, filename)

	// Create subdirectories if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
		return "", err
	}

	// Track for index
	relPath, _ := filepath.Rel(m.siteDir, filePath)
	m.pages = append(m.pages, PageInfo{
		URL:   pageURL,
		Title: title,
		File:  relPath,
	})

	return filePath, nil
}

// CreateIndex generates index.md file
func (m *Manager) CreateIndex() (string, error) {
	indexPath := filepath.Join(m.siteDir, "index.md")

	content := fmt.Sprintf(`---
site: %s
domain: %s
total_pages: %d
generated_at: %s
---

# Site Index: %s

**Source:** %s  
**Total Pages:** %d

## Pages

`,
		m.baseURL,
		m.domain,
		len(m.pages),
		time.Now().UTC().Format(time.RFC3339),
		m.domain,
		m.baseURL,
		len(m.pages))

	// Group by directory
	currentDir := ""
	for _, page := range m.pages {
		pageDir := filepath.Dir(page.File)
		if pageDir != currentDir && pageDir != "." {
			content += fmt.Sprintf("\n### %s/\n\n", pageDir)
			currentDir = pageDir
		}

		content += fmt.Sprintf("- [%s](%s)\n", page.Title, page.File)
		content += fmt.Sprintf("  - Source: %s\n", page.URL)
	}

	return indexPath, os.WriteFile(indexPath, []byte(content), 0644)
}

// SetMergeMode enables merge mode with the given batch size
func (m *Manager) SetMergeMode(enabled bool, batchSize int) {
	m.mergeMode = enabled
	m.mergeBatchSize = batchSize
	if m.mergeBatchSize <= 0 {
		m.mergeBatchSize = 50
	}
	m.mergeBuffer = nil
	m.mergePartNum = 0
	m.mergeTotalPages = 0
}

// IsMergeMode returns whether merge mode is enabled
func (m *Manager) IsMergeMode() bool {
	return m.mergeMode
}

// BufferPage adds a page to the merge buffer and flushes if the batch size is reached
func (m *Manager) BufferPage(pageURL, title, markdown string) error {
	m.mergeBuffer = append(m.mergeBuffer, MergedPage{
		URL:      pageURL,
		Title:    title,
		Markdown: markdown,
	})
	m.mergeTotalPages++

	if len(m.mergeBuffer) >= m.mergeBatchSize {
		return m.flushMergeBatch()
	}
	return nil
}

// Flush writes any remaining buffered pages to disk (call after crawl completes)
func (m *Manager) Flush() error {
	if !m.mergeMode || len(m.mergeBuffer) == 0 {
		return nil
	}
	return m.flushMergeBatch()
}

func (m *Manager) flushMergeBatch() error {
	if len(m.mergeBuffer) == 0 {
		return nil
	}

	m.mergePartNum++

	filename := fmt.Sprintf("output_part%d.md", m.mergePartNum)
	filePath := filepath.Join(m.outputDir, filename)

	if err := os.MkdirAll(m.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var buf strings.Builder

	// Write bundle header
	buf.WriteString(fmt.Sprintf("# Documentation Bundle: %s - Part %d\n\n", m.domain, m.mergePartNum))

	for i, page := range m.mergeBuffer {
		if i > 0 {
			buf.WriteString("\n---\n\n")
		}
		buf.WriteString(fmt.Sprintf("## %s\n\n", page.Title))
		buf.WriteString(fmt.Sprintf("**Source:** [%s](%s)\n\n", page.URL, page.URL))
		buf.WriteString(strings.TrimSpace(page.Markdown))
		buf.WriteString("\n")
	}

	if err := os.WriteFile(filePath, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("failed to write merge file %s: %w", filename, err)
	}

	// Clear buffer
	m.mergeBuffer = m.mergeBuffer[:0]

	return nil
}

func (m *Manager) generateFilename(pageURL, title string) string {
	// Try to use title first
	if title != "" && title != "Untitled" {
		slug := slugify(title)
		if slug != "" {
			return m.makeUnique(slug) + ".md"
		}
	}

	// Fallback to URL path
	parsed, err := url.Parse(pageURL)
	if err != nil {
		return m.makeUnique("page") + ".md"
	}

	path := strings.Trim(parsed.Path, "/")
	if path == "" {
		return "index.md"
	}

	// Use last segment
	segments := strings.Split(path, "/")
	lastSegment := segments[len(segments)-1]

	// Remove file extension if present
	if strings.Contains(lastSegment, ".") {
		lastSegment = strings.TrimSuffix(lastSegment, filepath.Ext(lastSegment))
	}

	slug := slugify(lastSegment)
	if slug == "" {
		slug = "page"
	}

	return m.makeUnique(slug) + ".md"
}

func (m *Manager) makeUnique(basename string) string {
	key := basename
	if count, exists := m.filenameCounts[key]; exists {
		m.filenameCounts[key] = count + 1
		return fmt.Sprintf("%s-%d", basename, count+1)
	}
	m.filenameCounts[key] = 1
	return basename
}

func slugify(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim
	slug = strings.Trim(slug, "-")

	// Truncate
	if len(slug) > 100 {
		slug = slug[:100]
		// Try to break at word boundary
		if idx := strings.LastIndex(slug, "-"); idx > 0 {
			slug = slug[:idx]
		}
	}

	return slug
}
