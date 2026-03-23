package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	quotes_scraping_scheduler "github.com/ozgeyilmazh/youtube-project/scrape-scheduler"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = quotes_scraping_scheduler.EnvironmentLocal
	}
	config := quotes_scraping_scheduler.NewConfig(environment)

	logLevel := slog.LevelInfo
	if environment == quotes_scraping_scheduler.EnvironmentProd {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))

	logger.Info("Starting scrape scheduler worker", "environment", environment)
	temporalClient, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Failed to dial temporal client", "error", err)
		os.Exit(1)
	}
	quotesScraperGatewayHost := config.GetScraperGatewayHost()
	quotesScraperGatewayClient := quotes_scraping_scheduler.NewQuotesScraperGatewayClient(quotesScraperGatewayHost, logger)
	scrapeActivity := quotes_scraping_scheduler.NewQuotesScrapeActivity(logger, quotesScraperGatewayClient)

	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 1 * time.Minute,
	}
	scraperWorkflow := quotes_scraping_scheduler.QuotesScrapeSchedulerWorkflow{
		Logger:          logger,
		ActivityOptions: activityOptions,
		ScrapeActivity:  scrapeActivity.BeginScrape,
	}

	worker := worker.New(temporalClient, config.GetSchedulerTaskQueue(), worker.Options{})
	worker.RegisterWorkflow(scraperWorkflow.Execute)
	worker.RegisterActivity(scrapeActivity.BeginScrape)
	err = worker.Start()
	if err != nil {
		logger.Error("Failed to start worker", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/healthz", quotes_scraping_scheduler.Healthz)
	http.HandleFunc("/livez", quotes_scraping_scheduler.Livez)
	http.HandleFunc("/readyz", quotes_scraping_scheduler.Readyz)
	go func() {
		logger.Info("Starting Healthz on port 8081")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			logger.Error("Failed to start healthz", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Worker running waiting for signals to stop ...")
	err = worker.Run(nil)
	if err != nil {
		logger.Error("Failed to run worker", "error", err)
		os.Exit(1)
	}
}
