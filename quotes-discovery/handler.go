package quotesdiscovery

//go:generate mockgen -destination=handler_mock.go -package=quotesdiscovery . HandlerWorkflow
import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	start, err := strconv.Atoi(c.Params("start"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	end, err := strconv.Atoi(c.Params("end"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for i := start; i <= end; i++ {
		workflowOptions := client.StartWorkflowOptions{
			ID:        uuid.New().String(),
			TaskQueue: h.taskQueue,
		}
		_, err := h.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, h.workflow.FetchQuotes, i)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.SendStatus(fiber.StatusAccepted)
}

// StartServer runs the HTTP app and a Temporal worker on the configured task queue.
// activity must be the same WorkflowActivity instance passed into NewWorkflow so activity
// names and behavior match what the workflow schedules.
func StartServer(workflow HandlerWorkflow, activity WorkflowActivity) (*fiber.App, worker.Worker) {

	config := NewConfig(EnvironmentLocal)
	app := fiber.New()
	taskQueue := config.GetTaskQueue()

	temporalClient, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
	})
	if err != nil {
		panic(err)
	}
	w := setupWorker(workflow, temporalClient, taskQueue, activity)

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

func setupWorker(workflow HandlerWorkflow, temporalClient client.Client, taskQueue string, activity WorkflowActivity) worker.Worker {
	worker := worker.New(temporalClient, taskQueue, worker.Options{})
	worker.RegisterWorkflow(workflow.FetchQuotes)
	worker.RegisterActivity(activity.FetchQuotes)
	worker.RegisterActivity(activity.BulkInsertData)
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
