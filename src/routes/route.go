package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"message-automation/src/handlers"
)

func InitializeRoutes(e *echo.Echo, h *handlers.MessageHandler) *echo.Echo {
	messageGroup := e.Group("/message")

	messageGroup.POST("/automation", h.HandleAutomation)
	messageGroup.GET("/messages", h.RetrieveMessages)
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
