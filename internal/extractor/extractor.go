package extractor

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// Extract extracts main content from HTML
func Extract(html, url string) (title, content, excerpt string, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", "", err
	}

	// Extract title
	title = doc.Find("title").First().Text()
	if title == "" {
		title = "Untitled"
	}

	// Try to get excerpt from meta tags
	excerpt = extractExcerpt(doc)

	// Remove unwanted elements
	doc.Find("script, style, nav, footer, header, aside").Remove()

	// Try to find main content
	var contentHTML string
	// Google Sites specific selector first (.tyJCtd), then standard semantic HTML
	mainSelectors := []string{".tyJCtd", "main", "article", "[role=main]", "#content", ".content"}
	for _, selector := range mainSelectors {
		selection := doc.Find(selector)
		if selection.Length() > 0 {
			// Extract ALL matching elements, not just the first one
			var parts []string
			selection.Each(func(i int, s *goquery.Selection) {
				if html, err := s.Html(); err == nil {
					parts = append(parts, html)
				}
			})
			contentHTML = strings.Join(parts, "\n")
			break
		}
	}

	// Fallback to body
	if contentHTML == "" {
		contentHTML, _ = doc.Find("body").Html()
	}

	return strings.TrimSpace(title), strings.TrimSpace(contentHTML), excerpt, nil
}

func extractExcerpt(doc *goquery.Document) string {
	// Try meta description
	if desc, exists := doc.Find("meta[name=description]").Attr("content"); exists && desc != "" {
		return strings.TrimSpace(desc)
	}

	// Try og:description
	if desc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists && desc != "" {
		return strings.TrimSpace(desc)
	}

	// Fallback to first paragraph
	if p := doc.Find("p").First(); p.Length() > 0 {
		text := p.Text()
		if len(text) > 50 {
			if len(text) > 200 {
				return text[:200] + "..."
			}
			return text
		}
	}

	return ""
}
