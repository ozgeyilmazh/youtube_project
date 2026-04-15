package quotesapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var (
	QUOTES_API_URL = "https://zenquotes.io/api/quotes?page="
)

func getLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

type Client struct {
	logger *slog.Logger
	client *http.Client
}

func NewClient(logger *slog.Logger) (*Client, error) {
	client := &Client{
		logger: logger,
		client: &http.Client{},
	}
	return client, nil
}

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type QuotesApiClientAdapter struct {
	client *Client
}

func NewQuotesApiClientAdapter(client *Client) *QuotesApiClientAdapter {
	return &QuotesApiClientAdapter{client: client}
}

func (a *QuotesApiClientAdapter) FetchQuotes(ctx context.Context, page int) ([]QuotesResponse, error) {
	a.client.logger.Info("Fetching quotes", "page", page)

	quotes, err := a.client.FetchQuotes(ctx, page)
	if err != nil {
		return nil, err
	}

	result := make([]QuotesResponse, len(quotes))
	for i, w := range quotes {
		result[i] = QuotesResponse{
			Quote:  w.Quote,
			Author: w.Author,
		}
	}
	return result, nil
}

func (c *Client) FetchQuotes(ctx context.Context, page int) ([]Quote, error) {
	c.logger.Info("Fetching quotes", "page", page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%d", QUOTES_API_URL, page), nil)
	if err != nil {
		return []Quote{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return []Quote{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return []Quote{}, fmt.Errorf("failed to fetch quotes: %s", resp.Status)
	}

	var quotes []Quote
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return []Quote{}, err
	}
	return quotes, nil
}
