package utils

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// HTMLMatcher represents the criteria used to match HTML elements for removal or processing.
// Tag specifies the HTML tag to match.
// Class specifies the CSS class to match.
// Attributes specifies a map of attribute key-value pairs to match.
type HTMLMatcher struct {
	Tag        string
	Class      string
	Attributes map[string]string
}

// CleanHTML removes unwanted HTML elements from the provided content based on the specified matchers.
func CleanHTML(content string, matchers []HTMLMatcher) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	// Remove unwanted elements
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		// Track nodes to remove
		var toRemove []*html.Node

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				for _, matcher := range matchers {
					if shouldRemove(c, matcher) {
						toRemove = append(toRemove, c)
					}
				}
			}
			traverse(c)
		}
		for _, node := range toRemove {
			n.RemoveChild(node)
		}
	}
	// Start with the root
	traverse(doc)

	// Render cleaned HTML
	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func shouldRemove(n *html.Node, matcher HTMLMatcher) bool {
	// Safety check - if no criteria, remove nothing.
	if matcher.Tag == "" && matcher.Class == "" && len(matcher.Attributes) == 0 {
		return false
	}

	// If only Tag is specified - with no class or attributes,
	// do a simple tag match
	if matcher.Tag != "" && matcher.Class == "" && len(matcher.Attributes) == 0 {
		return n.Data == matcher.Tag
	}

	// Match class
	if matcher.Class != "" {
		classFound := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, matcher.Class) {
				classFound = true
				break
			}
		}
		if !classFound {
			return false
		}
	}

	// Match attributes
	if len(matcher.Attributes) > 0 {
		// KV
		for expectedKey, expectedVal := range matcher.Attributes {
			attrFound := false
			for _, attr := range n.Attr {
				if attr.Key == expectedKey && attr.Val == expectedVal {
					attrFound = true
					break
				}
			}
			if !attrFound {
				return false
			}
		}
	}
	return true
}

// DefaultHTMLMatchers returns a slice of HTMLMatcher with default settings aimed at matching specific HTML elements.
func DefaultHTMLMatchers() []HTMLMatcher {
	return []HTMLMatcher{
		{
			Class: "gmail_quote",
		},
	}
}

func HTMLToMarkdown(html string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(html)
	if err != nil {
		return "", err
	}
	return markdown, nil
}
