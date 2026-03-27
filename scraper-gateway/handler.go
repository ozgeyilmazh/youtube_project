package scraper_gateway

//go:generate mockgen -destination=handler_mock.go -package=scraper_gateway . HandlerService
import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type HandlerService interface {
	TriggerScraping(ctx context.Context, start int, end int) error
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
	return nil
}
