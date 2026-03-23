package scrape_scheduler

import (
	"context"
	"fmt"
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
	url := fmt.Sprintf("%s/scrape/quotes/start/%d/end/%d", c.host, start, end)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to trigger scraping: %s", resp.Status)
	}
	return nil
}
