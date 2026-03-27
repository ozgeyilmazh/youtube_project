package scraper_gateway

import (
	"context"
	"net/http"
)

type QuotesDiscoveryConsumer struct {
	host   string
	client *http.Client
}

func NewQuotesDiscoveryConsumer(host string) *QuotesDiscoveryConsumer {
	return &QuotesDiscoveryConsumer{
		host:   host,
		client: &http.Client{},
	}
}

func (c *QuotesDiscoveryConsumer) TriggerScrape(ctx context.Context, start int, end int) error {
	return nil
}
