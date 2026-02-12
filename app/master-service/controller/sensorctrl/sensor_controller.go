package sensorctrl

import (
	"api/app/master-service/domain/sensordomain"
	"api/app/master-service/model"
	"api/utils/resp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorController struct {
	domain sensordomain.SensorDomain
}

func New(domain sensordomain.SensorDomain) *SensorController {
	return &SensorController{domain: domain}
}

func (ctl *SensorController) CreateForDevice(c *fiber.Ctx) error {
	deviceIDStr := c.Params("deviceId")
	deviceID, err := primitive.ObjectIDFromHex(deviceIDStr)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest("invalid device_id format"))
	}

	req := new(model.SensorCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	// Force relationship from route param
	req.DeviceID = deviceID

	sensor, err := ctl.domain.Create(c.Context(), req)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	return resp.Created(c, sensor)
}

func (ctl *SensorController) ListForDevice(c *fiber.Ctx) error {
	deviceID := c.Params("deviceId")
	limit := int64(100)
	if v := c.Query("limit"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil && i > 0 {
			limit = i
		}
	}
	sensors, err := ctl.domain.ListByDeviceID(c.Context(), deviceID, limit)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorInternal(err.Error()))
	}
	return resp.OK(c, sensors)
}

func (ctl *SensorController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	sensor, err := ctl.domain.GetByID(c.Context(), id)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorNotFound(err.Error()))
	}
	return resp.OK(c, sensor)
}

func (ctl *SensorController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(model.SensorUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	sensor, err := ctl.domain.Update(c.Context(), id, req)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	return resp.OK(c, sensor)
}

func (ctl *SensorController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctl.domain.Delete(c.Context(), id); err != nil {
		return resp.ErrorHandler(c, resp.ErrorNotFound(err.Error()))
	}
	return resp.OK(c)
}
