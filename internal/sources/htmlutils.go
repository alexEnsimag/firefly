package sources

import (
	"strings"

	"golang.org/x/net/html"
)

type Essay struct {
	Title, Description, Content string
}

func extractParagraphs(n *html.Node) string {
	var paragraphs []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			text := getText(c)
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		}
	}

	return strings.Join(paragraphs, "\n")
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var builder strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		builder.WriteString(getText(child))
	}

	return strings.TrimSpace(builder.String())
}

func getClassAttr(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return attr.Val
		}
	}
	return ""
}
