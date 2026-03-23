package scrape_scheduler

import (
	"log/slog"
	"testing"
)

func setupWorker(temporalClient client.Client, taskQueue string, quotesScrapeSchedulerWorkflow QuotesScrapeSchedulerWorkflow, logger *slog.Logger) {
	worker := worker.New(temporalClient, taskQueue, quotesScrapeSchedulerWorkflow)
	worker.RegisterWorkflow(quotesScrapeSchedulerWorkflow)
	worker.RegisterActivity(quotesScrapeSchedulerActivity)
	err := worker.Start()
	if err != nil {
		logger.Error("Failed to start worker", "error", err)
		panic(err)
	}
	return worker
}

func createTemporalClient(logger *slog.Logger) client.Client {

	environment := os.Getenv("ENV")
	if environment == "" {
		environment = ENVIRONMENT_LOCAL
	}

	config := NewConfig(environment)
	temporalClient, err := client.NewClient(client.Options{
		HostPort:  config.GetTemporalHostPort(),
		Logger:    logger,
	})
	if err != nil {
		logger.Error("Failed to create temporal client", "error", err)
		os.Exit(1)
	}
	return temporalClient
}