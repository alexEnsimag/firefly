package sources

import (
	"context"
	"firefly/alex/internal/httpclient"
	"testing"
)

func TestLoadEngadgeBlogPost(t *testing.T) {
	const (
		title       = "Sony and Yamaha are making a self-driving cart for theme parksIt's smarter, more comfortable and longer-lasting than its ancestors."
		description = "It's smarter, more comfortable and longer-lasting than its ancestors."
		content     = `Remember how we said Sony'sself-driving SC-1 conceptwould make for a great party bus? Apparently, Sony had the same idea. The company ispartneringwith Yamaha on the SC-1 Sociable Cart, an expansion of the concept designed for entertainment purposes like theme parks, golf courses and "commercial facilities." The new version seats five people instead of three (and in greater comfort), lasts longer through replaceable batteries and uses additional image sensors to improve its situational awareness.
As before, Sony feels the sensors eliminate the need for windows. A 49-inch 4K monitor on the inside provides a mixed reality view of the world, while four 55-inch 4K displays bombard passers-by with ads and other material. It will even use AI to optimize promos for outside people based on factors like age and gender -- not quiteMinority Reportlevels of eerily accurate ad targeting, but getting there.
The two companies expect to use the Sociable Cart for services in Japan sometime in fiscal 2019 (that is, before the end of March 2020). It won't, however, be available for sale. Not that you'd really want one given its glacial 11.8MPH top speed. This is strictly for fun on closed circuits, not your next pub crawl.`
	)

	client := httpclient.NewRateLimitedRetryClient(
		20, // requests per second
		1,  // no burst
		"MyApp/0.1-alpha",
	)

	e, err := LoadEngadgetBlogPost(context.Background(), client, "https://www.engadget.com/2019/08/25/sony-and-yamaha-sc-1-sociable-cart/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if e.Title != title {
		t.Errorf("wrong title")
	}
	if e.Description != description {
		t.Errorf("wrong description")
	}
	if e.Content != content {
		t.Errorf("wrong content")
	}
}
