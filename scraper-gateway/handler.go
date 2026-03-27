package scraper_gateway

//go:generate mockgen -destination=handler_mock.go -package=scraper_gateway . HandlerService
import (
	"context"
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
	start := ctx.Params("start")
	end := ctx.Params("end")
	go func() {
		startInt, _ := strconv.Atoi(start)
		endInt, _ := strconv.Atoi(end)
		err := h.service.TriggerScrape(ctx.Context(), startInt, endInt)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}()
	return ctx.SendStatus(fiber.StatusAccepted)
}
