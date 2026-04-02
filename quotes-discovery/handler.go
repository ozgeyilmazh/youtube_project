package quotesdiscovery

//go:generate mockgen -destination=handler_mock.go -package=quotesdiscovery . HandlerWorkflow
import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Quotes struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}
type HandlerWorkflow interface {
	FetchQuotes(ctx workflow.Context, page int) ([]Quotes, error)
}

type Handler struct {
	workflow          HandlerWorkflow
	temporalClient    client.Client
	logger            *slog.Logger
	config            *Config
	waitForCompletion bool
	taskQueue         string
}

func NewHandler(workflow HandlerWorkflow, temporalClient client.Client, logger *slog.Logger, config *Config, waitForCompletion bool, taskQueue string) *Handler {
	return &Handler{workflow: workflow, temporalClient: temporalClient, logger: logger, config: config, waitForCompletion: waitForCompletion, taskQueue: taskQueue}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/scrape/quotes/start/:start/end/:end", h.FetchQuotes)
}

func (h *Handler) FetchQuotes(c *fiber.Ctx) error {
	return nil
}

func StartServer(workflow HandlerWorkflow) (*fiber.App, worker.Worker) {

	config := NewConfig(EnvironmentLocal)
	app := fiber.New()
	taskQueue := config.GetTaskQueue()

	temporalClient, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
	})
	if err != nil {
		panic(err)
	}
	w := setupWorker(workflow, temporalClient, taskQueue)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	logger.Info("Starting quotes discovery server", "environment", EnvironmentLocal)
	handler := NewHandler(workflow, temporalClient, logger, config, false, taskQueue)
	handler.RegisterRoutes(app)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	return app, w
}

func setupWorker(workflow HandlerWorkflow, temporalClient client.Client, taskQueue string) worker.Worker {
	worker := worker.New(temporalClient, taskQueue, worker.Options{})
	worker.RegisterWorkflow(workflow.FetchQuotes)

	err := worker.Start()
	if err != nil {
		fmt.Println("Failed to start worker", err)
	}
	return worker
}

func WaitForServer(url string) error {
	maxRetries := 10
	client := http.Client{}

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/health", nil)

		if err != nil {
			cancel()
			return fmt.Errorf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		cancel()

		if err == nil {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Println("Failed to close response body", err)
				}
			}()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("failed to wait for completion")
}
