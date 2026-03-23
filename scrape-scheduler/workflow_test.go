package scrape_scheduler

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func setupWorker(temporalClient client.Client, taskQueue string, quotesScrapeSchedulerWorkflow QuotesScrapeSchedulerWorkflow, logger *slog.Logger) worker.Worker {
	worker := worker.New(temporalClient, taskQueue, worker.Options{})
	worker.RegisterWorkflow(quotesScrapeSchedulerWorkflow.Execute)
	worker.RegisterActivity(quotesScrapeSchedulerWorkflow.ScrapeActivity)
	err := worker.Start()
	if err != nil {
		logger.Error("Failed to start worker", "error", err)
		os.Exit(1)
	}
	return worker
}

func createTemporalClient(logger *slog.Logger) client.Client {

	environment := os.Getenv("ENV")
	if environment == "" {
		environment = EnvironmentLocal
	}

	config := NewConfig(environment)
	temporalClient, err := client.NewClient(client.Options{
		HostPort: config.GetTemporalHostPort(),
		Logger:   logger,
	})
	if err != nil {
		logger.Error("Failed to create temporal client", "error", err)
		os.Exit(1)
	}
	return temporalClient
}

func TestQuotesWorkflowSpec(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	temporalClient := createTemporalClient(logger)
	defer temporalClient.Close()

	taskQueue := "quotes-scrape-scheduler"
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: taskQueue,
	}

	Convey("Given a quotes scrape start at 1 end at 10 is scheduled", t, func() {

		start := 1
		end := 10
		isActivityCalled := false
		quotesScrapeSchedulerWorkflow := QuotesScrapeSchedulerWorkflow{
			Logger: logger,
			ActivityOptions: workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 1,
				},
			},
			ScrapeActivity: func(ctx context.Context, start int, end int) error {
				isActivityCalled = true
				return nil
			},
		}

		worker := setupWorker(temporalClient, taskQueue, quotesScrapeSchedulerWorkflow, logger)
		defer worker.Stop()

		Convey("When the workflow is started", func() {

			result, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, quotesScrapeSchedulerWorkflow.Execute, start, end)
			if err != nil {
				Convey("Then the workflow is failed", func() {
					So(err, ShouldNotBeNil)
				})
			}

			Convey("Then the workflow is scheduled", func() {
				quotesScrapeSchedulerResult := QuotesScrapeSchedulerResult{}
				err = result.Get(context.Background(), &quotesScrapeSchedulerResult)
				if err != nil {
					Convey("Then the workflow is failed", func() {
						So(err, ShouldNotBeNil)
					})
				}
				So(quotesScrapeSchedulerResult.Status, ShouldEqual, QuotesScraperStatusScheduled)
				So(isActivityCalled, ShouldBeTrue)
			})
		})
	})
}
