package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	quotes_discovery "github.com/ozgeyilmazh/youtube-project/quotes-discovery"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = quotes_discovery.EnvironmentLocal
	}

	logLevel := slog.LevelInfo
	if environment != quotes_discovery.EnvironmentProd {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	logger.Info("Starting quotes discovery api", "environment", environment)

	config := quotes_discovery.NewConfig(environment)
	temporalHostPort := config.GetTemporalHostPort()
	if temporalHostPort == "" {
		logger.Error("TEMPORAL_HOST_PORT is not set")
		os.Exit(1)
	}
	taskQueue := config.GetTaskQueue()
	if taskQueue == "" {
		logger.Error("SCHEDULER_TASK_QUEUE is not set")
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName: "Quotes Discovery",
	})
	app.Use(healthcheck.New())

	temporalClient, err := client.Dial(client.Options{
		HostPort: temporalHostPort,
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Failed to connect to Temporal", "error", err)
		os.Exit(1)
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}

	// Setup database
	db, tableName, err := quotes_discovery.SetupClickhouseRepository(config.GetClickhouseTableName())
	if err != nil {
		logger.Error("Unable to setup clickhouse repository", "error", err)
		os.Exit(1)
	}
	repository := quotes_discovery.NewClickhouseRepository(db, tableName)

	quotesClient := quotes_discovery.NewQuotesAPIClient(config.GetQuotesClientHost(), logger)
	activity := quotes_discovery.NewActivity(quotesClient, repository)
	workflow := quotes_discovery.NewWorkflow(activity, activityOptions, config)

	handler := quotes_discovery.NewHandler(workflow, temporalClient, logger, config, false, taskQueue)
	handler.RegisterRoutes(app)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", config.GetPort())); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
