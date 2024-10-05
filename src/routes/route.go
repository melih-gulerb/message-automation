package routes

import (
	"github.com/labstack/echo/v4"
	"message-automation/src/handlers"
)

func InitializeRoutes(e *echo.Echo, h *handlers.MessageHandler) *echo.Echo {
	e.POST("/automation", h.HandleAutomation)
	e.GET("/retrieve", h.RetrieveMessages)
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running")
	})

	return e
}
