package sensorctrl

import (
	"api/app/master-service/domain/sensordomain"
	"api/app/master-service/model"
	"api/lib"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

type SensorController struct {
	sensorDomain sensordomain.SensorDomain
}

func New(sensorDomain sensordomain.SensorDomain) *SensorController {
	return &SensorController{sensorDomain: sensorDomain}
}

// GetAll gets all sensors for a device
// @Summary Get all sensors for a device
// @Description Get all sensors for a device
// @Tags Sensor
// @Accept json
// @Produce json
// @Param device_id path string true "Device ID"
// @Success 200 {object} resp.Response{data=[]model.Sensor} "List of Sensors"
// @Router /devices/{device_id}/sensors [get]
// @Security TokenKey
func (c *SensorController) GetAll(ctx *fiber.Ctx) error {
	sensors, err := c.sensorDomain.GetAll(ctx.Context(), lib.StrToUUID(ctx.Params("device_id")))
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, sensors)
}

// Create creates a new sensor for a device
// @Summary Create a new sensor for a device
// @Description Create a new sensor for a device
// @Tags Sensor
// @Accept json
// @Produce json
// @Param device_id path string true "Device ID"
// @Param sensor body model.SensorAPI true "Body Request"
// @Success 201 {object} resp.Response{data=model.Sensor} "Sensor data"
// @Router /devices/{device_id}/sensors [post]
// @Security TokenKey
func (c *SensorController) Create(ctx *fiber.Ctx) error {
	deviceID := lib.StrToUUID(ctx.Params("device_id"))

	api := new(model.SensorAPI)
	if err := ctx.BodyParser(api); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	// Force relationship from route param
	api.DeviceID = deviceID

	sensor, err := c.sensorDomain.Create(ctx.Context(), api)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.Created(ctx, sensor)
}

// GetByID gets a sensor by id
// @Summary Get a sensor by id
// @Description Get a sensor by id
// @Tags Sensor
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Param device_id path string true "Device ID"
// @Success 200 {object} resp.Response{data=model.Sensor} "Sensor data"
// @Router /devices/{device_id}/sensors/{id} [get]
// @Security TokenKey
func (c *SensorController) GetByID(ctx *fiber.Ctx) error {
	sensor, err := c.sensorDomain.GetByID(ctx.Context(), lib.StrToUUID(ctx.Params("device_id")), lib.ParamsUUID(ctx))
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, sensor)
}

// Update updates a sensor
// @Summary Update a sensor
// @Description Update a sensor
// @Tags Sensor
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Param device_id path string true "Device ID"
// @Param sensor body model.SensorUpdateRequest true "Body Request"
// @Success 200 {object} resp.Response{data=model.Sensor} "Sensor data"
// @Router /devices/{device_id}/sensors/{id} [put]
// @Security TokenKey
func (c *SensorController) Update(ctx *fiber.Ctx) error {
	api := new(model.SensorUpdateRequest)
	if err := ctx.BodyParser(api); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	sensor, err := c.sensorDomain.Update(ctx.Context(), lib.StrToUUID(ctx.Params("device_id")), lib.ParamsUUID(ctx), api)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, sensor)
}

// Delete deletes a sensor
// @Summary Delete a sensor
// @Description Delete a sensor
// @Tags Sensor
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Param device_id path string true "Device ID"
// @Success 200 {object} resp.Response "Success"
// @Router /devices/{device_id}/sensors/{id} [delete]
// @Security TokenKey
func (c *SensorController) Delete(ctx *fiber.Ctx) error {
	if err := c.sensorDomain.Delete(ctx.Context(), lib.StrToUUID(ctx.Params("device_id")), lib.ParamsUUID(ctx)); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx)
}
