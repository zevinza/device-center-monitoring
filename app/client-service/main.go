package main

import (
	"api/app/client-service/config"
	"api/app/client-service/model"
	"api/utils/env"
	"log"
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
			"service": viper.GetString("NAME"),
			"status":  "running",
			"port":    viper.GetInt("PORT"),
		})
	})

	// Receive sensor data endpoint
	app.Post("/receive", receiveSensorData)

	port := viper.GetString("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Client Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func receiveSensorData(c *fiber.Ctx) error {
	req := new(model.SensorDataRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	// Log the received data with device and sensor details
	log.Printf("[CLIENT] Received sensor data:")
	log.Printf("  ReadingID: %s", req.ReadingID)
	log.Printf("  Device: %s (Code: %s, ID: %s)", req.Device.Name, req.Device.DeviceCode, req.Device.ID)
	if req.Sensor != nil {
		log.Printf("  Sensor: %s (ID: %s, Type: %s)", req.Sensor.Name, req.Sensor.ID, req.Sensor.Type)
	}
	log.Printf("  Value: %.2f %s", req.Value, req.Unit)
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
