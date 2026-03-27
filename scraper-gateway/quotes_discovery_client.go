package scraper_gateway

import (
	"context"
	"fmt"
	"net/http"
)

type QuotesDiscoveryClient struct {
	host   string
	client *http.Client
}

func NewQuotesDiscoveryClient(host string) *QuotesDiscoveryClient {
	return &QuotesDiscoveryClient{
		host:   host,
		client: &http.Client{},
	}
}

func (c *QuotesDiscoveryClient) TriggerScrape(ctx context.Context, start int, end int) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/scrape/quotes/start/%d/end/%d", c.host, start, end), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("expected status %d but got %s", http.StatusAccepted, resp.Status)
	}
	return nil
}
