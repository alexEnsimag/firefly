package sources

import (
	"context"
	"firefly/alex/internal/httpclient"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/html"
)

const engadgetBaseURL = "https://www.engadget.com"

func LoadEngadgetBlogPost(ctx context.Context, httpClient *httpclient.RateLimitedRetryClient, u string) (*Essay, error) {
	if !strings.HasPrefix(u, engadgetBaseURL) {
		return nil, fmt.Errorf("supports only essays from %s", engadgetBaseURL)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing html body: %w", err)
	}

	var e Essay
	parseEngadgetBlogPost(doc, &e)
	if e.Title == "" || e.Description == "" || e.Content == "" {
		return nil, fmt.Errorf("failed parsing essay")
	}

	return &e, nil
}

func parseEngadgetBlogPost(n *html.Node, e *Essay) {
	if n.Type == html.ElementNode {
		class := getClassAttr(n)

		switch {
		case strings.Contains(class, "caas-title"):
			e.Title = getText(n)

		case strings.Contains(class, "caas-subhead"):
			e.Description = getText(n)

		case strings.Contains(class, "caas-body"):
			e.Content = extractParagraphs(n)
			return // stop searching after finding content
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		parseEngadgetBlogPost(child, e)
	}
}
