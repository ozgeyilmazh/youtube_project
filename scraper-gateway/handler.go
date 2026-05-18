package scraper_gateway

//go:generate mockgen -destination=handler_mock.go -package=scraper_gateway . HandlerService
import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HandlerService interface {
	TriggerScrape(ctx context.Context, start int, end int) error
}

type Handler struct {
	service HandlerService
}

func NewHandler(service HandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/scrape/quotes/start/:start/end/:end", h.TriggerScraping)
}

func (h *Handler) TriggerScraping(ctx *fiber.Ctx) error {
	startInt, err := strconv.Atoi(ctx.Params("start"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid start",
		})
	}
	endInt, err := strconv.Atoi(ctx.Params("end"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid end",
		})
	}

	go func() {
		if err := h.service.TriggerScrape(context.Background(), startInt, endInt); err != nil {
			slog.Error("trigger scrape failed", "start", startInt, "end", endInt, "error", err)
		}
	}()
	return ctx.SendStatus(fiber.StatusAccepted)
}
