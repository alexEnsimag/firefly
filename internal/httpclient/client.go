package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
)

type RateLimitedRetryClient struct {
	client  *retryablehttp.Client
	limiter *rate.Limiter
}

func NewRateLimitedRetryClient(rps float64, burst int, userAgent string) *RateLimitedRetryClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Logger = nil

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return true, nil
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			return true, nil
		}
		return false, nil
	}

	retryClient.HTTPClient.Transport = &rateLimitedTransport{
		limiter:   rate.NewLimiter(rate.Limit(rps), burst),
		base:      http.DefaultTransport,
		userAgent: userAgent,
	}

	return &RateLimitedRetryClient{
		client:  retryClient,
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

func (c *RateLimitedRetryClient) Get(url string) (*http.Response, error) {
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *RateLimitedRetryClient) Do(req *retryablehttp.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// --- Internal Transport Wrapper --- //

type rateLimitedTransport struct {
	limiter   *rate.Limiter
	base      http.RoundTripper
	userAgent string
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	return t.base.RoundTrip(req)
}
