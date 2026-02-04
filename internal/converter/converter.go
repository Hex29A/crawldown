package converter

import (
	"fmt"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// ToMarkdown converts HTML to markdown
func ToMarkdown(html string) (string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", err
	}

	return cleanMarkdown(markdown), nil
}

// FormatWithFrontmatter adds frontmatter to markdown
func FormatWithFrontmatter(title, url, markdown, excerpt string) string {
	frontmatter := fmt.Sprintf(`---
source_url: %s
title: %s
crawled_at: %s`,
		url,
		title,
		time.Now().UTC().Format(time.RFC3339))

	if excerpt != "" {
		// Escape quotes in excerpt
		excerpt = strings.ReplaceAll(excerpt, `"`, `\"`)
		excerpt = strings.ReplaceAll(excerpt, "\n", " ")
		frontmatter += fmt.Sprintf("\nexcerpt: \"%s\"", excerpt)
	}

	frontmatter += "\n---\n\n"

	// Ensure title as H1
	content := strings.TrimSpace(markdown)
	if !strings.HasPrefix(content, fmt.Sprintf("# %s", title)) {
		content = fmt.Sprintf("# %s\n\n%s", title, content)
	}

	return frontmatter + content
}

func cleanMarkdown(markdown string) string {
	// Remove excessive blank lines
	lines := strings.Split(markdown, "\n")
	var cleaned []string
	prevBlank := false

	for _, line := range lines {
		isBlank := strings.TrimSpace(line) == ""
		if isBlank {
			if !prevBlank {
				cleaned = append(cleaned, "")
			}
			prevBlank = true
		} else {
			cleaned = append(cleaned, line)
			prevBlank = false
		}
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
