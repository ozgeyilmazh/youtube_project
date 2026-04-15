package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	quotesapi "github.com/ozgeyilmazh/youtube-project/quotes-api"
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}

	logLevel := slog.LevelInfo
	if environment == "prod" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))

	logger.Info("Starting quotes API", "environment", environment)

	config := quotesapi.NewConfig(environment)
	app := fiber.New(fiber.Config{
		AppName: "quotes-api",
	})

	quotesApiClient, err := quotesapi.NewClient(logger)
	if err != nil {
		logger.Error("Failed to create quotes API client", "error", err)
		os.Exit(1)
	}
	quotesApiClientAdapter := quotesapi.NewQuotesApiClientAdapter(quotesApiClient)
	service := quotesapi.NewService(quotesApiClientAdapter)
	handler := quotesapi.NewHandler(service)

	app.Use(healthcheck.New())
	handler.RegisterRoutes(app)
	err = app.Listen(fmt.Sprintf("0.0.0.0:%s", config.GetPort()))
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
