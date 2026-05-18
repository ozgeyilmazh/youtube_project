package main

import (
	"log/slog"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	quotes_discovery "github.com/ozgeyilmazh/youtube-project/quotes-discovery"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = quotes_discovery.EnvironmentLocal
	}
	config := quotes_discovery.NewConfig(environment)

	logLevel := slog.LevelInfo
	if environment != quotes_discovery.EnvironmentProd {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	logger.Info("Starting suspect discovery worker...")

	temporalClient, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Unable to create Temporal Client", "error", err)
		os.Exit(1)
	}

	db, tableName, err := quotes_discovery.SetupClickhouseRepository(config.GetClickhouseTableName())
	if err != nil {
		logger.Error("Unable to setup database", "error", err)
		os.Exit(1)
	}
	repository := quotes_discovery.NewClickhouseRepository(db, tableName)

	quotesClient := quotes_discovery.NewQuotesAPIClient(config.GetQuotesClientHost(), logger)
	activity := quotes_discovery.NewActivity(quotesClient, repository)

	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 3 * time.Minute,
	}

	workflowImpl := quotes_discovery.NewWorkflow(activity, activityOptions, config)

	taskQueue := config.GetTaskQueue()
	if taskQueue == "" {
		logger.Error("SCHEDULER_TASK_QUEUE is not set")
		os.Exit(1)
	}
	logger.Info("Worker will listen on queue", "queue", taskQueue)
	worker := worker.New(temporalClient, taskQueue, worker.Options{})
	worker.RegisterWorkflow(workflowImpl.FetchQuotes)
	worker.RegisterActivity(activity.FetchQuotes)
	worker.RegisterActivity(activity.BulkInsertData)

	app := fiber.New()
	app.Use(healthcheck.New())
	go func() {
		logger.Info("Starting Healthz on port 3000")
		err := app.Listen("0.0.0.0:3000")
		if err != nil {
			logger.Error("Unable to start Healthz", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Worker running, waiting for signals to stop...")
	err = worker.Run(nil)
	if err != nil {
		logger.Error("Worker stopped with error", "error", err)
		os.Exit(1)
	}
}
