package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	scraper_gateway "github.com/ozgeyilmazh/youtube-project/scraper-gateway"
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

	logger.Info("Starting scraper gateway", "environment", environment)

	config := scraper_gateway.NewConfig(environment)
	app := fiber.New(fiber.Config{
		AppName: "scraper-gateway",
	})

	quotesDiscoveryHost := config.GetQuotesDiscoveryHost()
	quotesDiscoveryClient := scraper_gateway.NewQuotesDiscoveryClient(quotesDiscoveryHost)
	service := scraper_gateway.NewService(quotesDiscoveryClient)
	handler := scraper_gateway.NewHandler(service)

	app.Use(healthcheck.New())

	handler.RegisterRoutes(app)

	err := app.Listen(fmt.Sprintf("0.0.0.0:%s", config.GetPort()))
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
