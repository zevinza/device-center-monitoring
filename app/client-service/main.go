package main

import (
	"api/app/client-service/config"
	"api/app/client-service/model"
	"api/utils/env"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func init() {
	env.Load(config.Environment)
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Use(cors.New())
	app.Use(recover.New())

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":     viper.GetString("NAME"),
			"description": viper.GetString("DESCRIPTION"),
			"status":      "running",
		})
	})

	// Receive sensor data endpoint
	app.Post("/receive", receiveSensorData)

	port := viper.GetString("PORT")

	log.Printf("Client Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func receiveSensorData(c *fiber.Ctx) error {
	req := new(model.SensorReceiveRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if rand.Intn(100) < 30 {
		log.Printf("[CLIENT] Failed to receive sensor data: %v", req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": "Failed to receive sensor data",
		})
	}

	log.Printf("[CLIENT] Received sensor data:")
	log.Printf("  ReadingID: %s", req.ID)
	log.Printf("  Device: %s (Code: %s, ID: %s)", req.Device.Name, req.Device.Code, req.Device.ID)
	if req.Sensor != nil {
		log.Printf("  Sensor: %s (ID: %s, Unit: %s)", req.Sensor.Name, req.Sensor.ID, req.Sensor.Unit)
	}
	log.Printf("  Value: %s %s", req.Value, req.Sensor.Unit)
	log.Printf("  Timestamp: %s", req.Timestamp.Format(time.RFC3339))

	// In a real scenario, you would save this to your own storage here
	// For now, we just log it as specified in the plan

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"message":     "Sensor data received",
		"data":        req,
		"received_at": time.Now().UTC().Format(time.RFC3339),
	})
}
