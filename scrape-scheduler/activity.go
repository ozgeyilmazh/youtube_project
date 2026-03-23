package scrape_scheduler

//go:generate mockgen -destination=activity_mock.go -package=scrape_scheduler . Client

import (
	"context"
	"log/slog"
)

type (
	Client interface {
		TriggerScraping(ctx context.Context, start int, end int) error
	}
	QuotesScrapeActivity struct {
		logger *slog.Logger
		client Client
	}
)

func NewQuotesScrapeActivity(logger *slog.Logger, client Client) *QuotesScrapeActivity {
	return &QuotesScrapeActivity{logger: logger, client: client}
}

func (a *QuotesScrapeActivity) BeginScrape(ctx context.Context, start int, end int) error {
	a.logger.Info("Beginning quotes scrape", "start", start, "end", end)
	return a.client.TriggerScraping(ctx, start, end)
}
