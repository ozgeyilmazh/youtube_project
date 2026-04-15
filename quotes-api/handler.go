package quotesapi

//go:generate mockgen -destination=handler_mock.go -package=quotesapi . HandlerService

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type QuotesResponse struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type HandlerService interface {
	FetchQuotes(ctx context.Context, page int) ([]QuotesResponse, error)
}

type Handler struct {
	service HandlerService
}

func NewHandler(service HandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/fetch/quotes", h.FetchQuotes)
}

func (h *Handler) FetchQuotes(c *fiber.Ctx) error {
	return nil
}
