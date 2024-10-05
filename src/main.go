package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"message-automation/src/clients"
	"message-automation/src/configs"
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
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	messageRepository := repositories.NewMessageRepository(db)
	webhookClient := clients.NewWebhookClient(cfg.WebhookBaseURL, cfg.WebhookToken)
	redisClient := clients.NewRedisClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB, cfg.RedisTimeout)
	messageService := services.NewMessageService(messageRepository, webhookClient, redisClient)
	messageHandler := handlers.NewMessageHandler(messageService)

	routes.InitializeRoutes(e, messageHandler)

	// Execute all unsent messages with the project deployment
	messageService.ExecuteAutomationForProjectDeployment()

	// Execute message automation async with 2 seconds period
	go messageService.ExecuteAutomation(services.IsActiveChannel, 2)

	e.Logger.Fatal(e.Start(":3030"))
}
