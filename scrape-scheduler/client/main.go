package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	quotes_scraping_scheduler "github.com/ozgeyilmazh/youtube-project/scrape-scheduler"
	"go.temporal.io/sdk/client"
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

	logger.Info("Starting scrape scheduler client", "environment", environment)

	schedulerCron := config.GetSchedulerCron()
	if schedulerCron == "" {
		logger.Error("Scheduler cron is not set")
		os.Exit(1)
	}

	quotesScraperGatewayHost := config.GetScraperGatewayHost()
	quotesScraperGatewayClient := quotes_scraping_scheduler.NewQuotesScraperGatewayClient(quotesScraperGatewayHost, logger)
	scrapeActivity := quotes_scraping_scheduler.NewQuotesScrapeActivity(logger, quotesScraperGatewayClient)

	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	scraperWorkflow := quotes_scraping_scheduler.QuotesScrapeSchedulerWorkflow{
		Logger:          logger,
		ActivityOptions: activityOptions,
		ScrapeActivity:  scrapeActivity.BeginScrape,
	}

	workflowOptions := client.StartWorkflowOptions{
		CronSchedule: schedulerCron,
		ID:           uuid.New().String(),
		TaskQueue:    config.GetSchedulerTaskQueue(),
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Failed to dial temporal client", "error", err)
		os.Exit(1)
	}

	start := 1
	end := 10

	quotesScrapingSchedulerWorkflowID, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, scraperWorkflow.Execute, start, end)
	if err != nil {
		logger.Error("Failed to execute workflow", "error", err)
		os.Exit(1)
	}

	logger.Info("Workflow started", "workflowID", quotesScrapingSchedulerWorkflowID)

	http.HandleFunc("/healthz", quotes_scraping_scheduler.Healthz)
	go func() {
		logger.Info("Starting Healthz on port 8081")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			logger.Error("Failed to start healthz", "error", err)
			os.Exit(1)
		}
	}()
	logger.Info("Started workflow", "workflowID", quotesScrapingSchedulerWorkflowID)

	for {
		time.Sleep(100 * time.Second)
	}
}
