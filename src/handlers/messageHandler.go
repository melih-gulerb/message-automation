package handlers

import (
	"github.com/labstack/echo/v4"
	"message-automation/src/models"
	"message-automation/src/models/base"
	"message-automation/src/services"
)

type MessageHandler struct {
	MessageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	h := MessageHandler{MessageService: messageService}

	return &h
}

func (h *MessageHandler) HandleAutomation(c echo.Context) error {
	isActive := c.QueryParam("isActive")
	h.MessageService.HandleAutomation(isActive)

	return c.JSON(200, map[string]string{"success": "true"})
}

func (h *MessageHandler) RetrieveMessages(c echo.Context) error {
	limit := c.QueryParam("limit")

	messages := h.MessageService.RetrieveSentMessages(limit)

	return c.JSON(200, base.Response[[]models.Message]{
		Success: true,
		Data:    messages,
	})
}
