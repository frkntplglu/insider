package handler

import (
	"context"

	"github.com/frkntplglu/insider/internal/models"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/gofiber/fiber/v2"
)

type messageService interface {
	GetAllSentMessages(ctx context.Context) ([]models.MessageSentItem, error)
	StartAutoSending() error
	StopAutoSending() error
}

type MessageHandler struct {
	messageService messageService
}

func NewMessageHandler(messageService messageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) GetAllSentMessages(c *fiber.Ctx) error {
	messages, err := h.messageService.GetAllSentMessages(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(FailureResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Data:    messages,
	})
}

func (h *MessageHandler) StartAutoSending(c *fiber.Ctx) error {
	err := h.messageService.StartAutoSending()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(FailureResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Data:    "OK",
	})
}

func (h *MessageHandler) StopAutoSending(c *fiber.Ctx) error {
	err := h.messageService.StopAutoSending()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(FailureResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
		Data:    "OK",
	})
}

func (h *MessageHandler) SetRoutes(app *fiber.App) {
	messageGroup := app.Group("/messages")
	messageGroup.Get("/", h.GetAllSentMessages)
	messageGroup.Get("/start", h.StartAutoSending)
	messageGroup.Get("/stop", h.StopAutoSending)
}

var MessageSwaggerEndpoints = []*endpoint.EndPoint{
	endpoint.New(
		endpoint.GET,
		"/messages",
		endpoint.WithTags("Messages"),
		endpoint.WithSuccessfulReturns([]response.Response{response.New([]models.MessageSentItem{}, "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(FailureResponse{}, "Bad Request", "400")}),
		endpoint.WithDescription("It returns sent messages from Redis cache"),
	),
	endpoint.New(
		endpoint.GET,
		"/messages/start",
		endpoint.WithTags("Messages"),
		endpoint.WithSuccessfulReturns([]response.Response{response.New("OK", "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(FailureResponse{}, "Bad Request", "400")}),
		endpoint.WithDescription("It starts autosending mechanism if it is stopped"),
	),
	endpoint.New(
		endpoint.GET,
		"/messages/stop",
		endpoint.WithTags("Messages"),
		endpoint.WithSuccessfulReturns([]response.Response{response.New("OK", "OK", "200")}),
		endpoint.WithErrors([]response.Response{response.New(FailureResponse{}, "Bad Request", "400")}),
		endpoint.WithDescription("It stops autosending mechanism if it is started"),
	),
}
