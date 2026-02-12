package main

import (
	"api/app/master-service/config"
	"api/app/master-service/initialization"
	"api/app/master-service/routes"
	"api/middleware"
	"api/services/cache"
	"api/services/database"
	"api/utils/env"
	"api/utils/logging"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	env.Load(config.Environment)
	cache.Connect()
	database.ConnectMongo()
}

// @title API
// @version 1.0.0
// @description API Documentation
// @contact.name Armada Muhammad Siswanto
// @contact.email armadamuhammads@gmail.com
// @host localhost:8000
// @schemes http
// @BasePath /api/v1/master
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey TokenKey
// @in header
// @name X-API-Key
func main() {
	logging.Init()

	app := fiber.New(fiber.Config{
		Prefork: viper.GetString("PREFORK") == "true",
	})

	app.Use(middleware.RequestLog())
	app.Use(middleware.RecoverSlog())

	module := initialization.Init()

	routes.Handle(app, module)

	// Start queue consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if module.Service.QueueConsumerService != nil {
		go module.Service.QueueConsumerService.Start(ctx)
		log.Println("Queue consumer started")
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down...")
		cancel()
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	log.Fatal(app.Listen(":" + viper.GetString("PORT")))
}
