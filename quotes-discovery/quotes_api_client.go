package quotesdiscovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type QuotesAPIClient struct {
	host   string
	client *http.Client
	logger *slog.Logger
}

func NewQuotesAPIClient(host string, logger *slog.Logger) *QuotesAPIClient {
	return &QuotesAPIClient{
		host:   host,
		logger: logger,
		client: &http.Client{},
	}
}

func (c *QuotesAPIClient) FetchQuotes(ctx context.Context, page int) ([]Quotes, error) {
	c.logger.Info("Fetching quotes", "page", page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/fetch/quotes?page=%d", c.host, page), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch quotes: %s", resp.Status)
	}
	var quotes []Quotes
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}
