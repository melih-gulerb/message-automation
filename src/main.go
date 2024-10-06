package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"message-automation/src/clients"
	"message-automation/src/configs"
	_ "message-automation/src/docs"
	"message-automation/src/handlers"
	"message-automation/src/middlewares"
	"message-automation/src/repositories"
	"message-automation/src/routes"
	"message-automation/src/services"
)

func main() {
	cfg := configs.SetConfig()
	db := configs.InitDB(cfg.MSSQLConnectionString)
	e := echo.New()
	e.Use(middlewares.RecoverMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	messageRepository := repositories.NewMessageRepository(db)
	webhookClient := clients.NewWebhookClient(cfg.WebhookBaseURL, cfg.WebhookToken)
	redisClient := clients.NewRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB, cfg.RedisTimeout)

	messageService := services.NewMessageService(messageRepository, webhookClient, redisClient, cfg.MessagePerExecution, cfg.ExecutionPeriod)
	messageHandler := handlers.NewMessageHandler(messageService)

	routes.InitializeRoutes(e, messageHandler)

	// Execute all unsent messages with the project deployment
	messageService.ExecuteAutomationForProjectDeployment()

	// Execute message automation asynchronously with a specific period
	go messageService.ExecuteAutomation()

	e.Logger.Fatal(e.Start(":3030"))
}
