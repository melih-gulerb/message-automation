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

// HandleAutomation godoc
// @Summary Automates message sending operations based on isActive status
// @Description Trigger message automation based on the isActive query parameter
// @Tags messages
// @Accept  json
// @Produce  json
// @Param   isActive   query   string   true  "Automation activation status"
// @Success 200 {object} base.Response[string]
// @Failure 400 {object} base.Response[base.Error]
// @Failure 500 {object} base.Response[base.Error]
// @Router /message/automation [post]
func (h *MessageHandler) HandleAutomation(c echo.Context) error {
	isActive := c.QueryParam("isActive")
	h.MessageService.HandleAutomation(isActive)

	return c.JSON(200, base.Response[string]{
		Success: true,
		Data:    "succeed",
	})
}

// RetrieveMessages godoc
// @Summary Retrieve sent messages
// @Description Retrieve a list of sent messages, limited by the query parameter and messageId
// @Tags messages
// @Accept  json
// @Produce  json
// @Param   limit   query   string   false  "Limit number of messages"
// @Param   messageId   query   string   false  "Filter by Id of the message"
// @Success 200 {object} base.Response[[]models.Message]
// @Failure 400 {object} base.Response[base.Error]
// @Failure 500 {object} base.Response[base.Error]
// @Router /message/messages [get]
func (h *MessageHandler) RetrieveMessages(c echo.Context) error {
	limit := c.QueryParam("limit")
	messageId := c.QueryParam("messageId")

	messages := h.MessageService.RetrieveSentMessages(limit, messageId)

	return c.JSON(200, base.Response[[]models.Message]{
		Success: true,
		Data:    messages,
	})
}
