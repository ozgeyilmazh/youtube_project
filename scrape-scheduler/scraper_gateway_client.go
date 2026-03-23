package scrape_scheduler

import (
	"context"
	"log/slog"
	"net/http"
)

type QuotesScraperGatewayClient struct {
	host   string
	logger *slog.Logger
	client *http.Client
}

func NewQuotesScraperGatewayClient(host string, logger *slog.Logger) *QuotesScraperGatewayClient {
	return &QuotesScraperGatewayClient{
		host:   host,
		logger: logger,
		client: &http.Client{},
	}
}

func (c *QuotesScraperGatewayClient) TriggerScraping(ctx context.Context, start int, end int) error {
	c.logger.Info("Triggering scraping", "start", start, "end", end)
	return nil
}
