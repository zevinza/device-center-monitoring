package routes

import (
	"api/app/master-service/initialization"
	"api/lib"
	"api/middleware"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func Handle(app *fiber.App, module initialization.Module) {
	app.Use(cors.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			lib.PrintStackTrace(e)
		},
	}))

	api := app.Group(viper.GetString("ENDPOINT"))

	api.Get("/", func(c *fiber.Ctx) error {
		return resp.OK(c, viper.GetString("NAME"))
	})

	controllers := module.Controller

	// Device
	deviceController := controllers.Device
	deviceAPI := api.Group("/devices").Use(middleware.SecretKeyAuthentication())
	deviceAPI.Post("/", deviceController.Create)
	deviceAPI.Get("/", deviceController.GetAll)
	deviceAPI.Get("/:id", deviceController.GetByID)
	deviceAPI.Put("/:id", deviceController.Update)
	deviceAPI.Delete("/:id", deviceController.Delete)

	// Sensor
	sensorController := controllers.Sensor
	sensorAPI := api.Group("/devices/:device_id/sensors").Use(middleware.SecretKeyAuthentication())
	sensorAPI.Post("/", sensorController.CreateForDevice)
	sensorAPI.Get("/", sensorController.ListForDevice)
	sensorAPI.Get("/:id", sensorController.GetByID)
	sensorAPI.Put("/:id", sensorController.Update)
	sensorAPI.Delete("/:id", sensorController.Delete)

	// Sensor Ingest
	sensorIngestController := controllers.SensorIngest
	sensorIngestAPI := api.Group("/sensors").Use(middleware.SecretKeyAuthentication())
	sensorIngestAPI.Post("/", sensorIngestController.Ingest)
}
