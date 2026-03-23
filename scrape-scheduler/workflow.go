package scrape_scheduler

import (
	"context"
	"log/slog"

	"go.temporal.io/sdk/workflow"
)

type QuotesScrapeSchedulerState string

type QuotesScrapeSchedulerResult struct {
	Status QuotesScrapeSchedulerState
}

const (
	QuotesScraperStatusScheduled QuotesScrapeSchedulerState = "scheduled"
	QuotesScraperStatusFailed    QuotesScrapeSchedulerState = "failed"
)

type QuotesScrapeSchedulerWorkflow struct {
	Logger          *slog.Logger
	ActivityOptions workflow.ActivityOptions
	ScrapeActivity  func(ctx context.Context, start int, end int) error
}

func (w *QuotesScrapeSchedulerWorkflow) Execute(ctx workflow.Context, start int, end int) (QuotesScrapeSchedulerResult, error) {
	w.Logger.Info("Executing quotes scrape scheduler workflow", "start", start, "end", end)

	ctx = workflow.WithActivityOptions(ctx, w.ActivityOptions)
	err := workflow.ExecuteActivity(ctx, w.ScrapeActivity, start, end).Get(ctx, nil)
	if err != nil {
		w.Logger.Error("Failed to execute quotes scrape activity", "error", err)
		return QuotesScrapeSchedulerResult{Status: QuotesScraperStatusFailed}, err
	}

	w.Logger.Info("Quotes scrape scheduler workflow executed successfully", "start", start, "end", end)
	return QuotesScrapeSchedulerResult{Status: QuotesScraperStatusScheduled}, nil
}
