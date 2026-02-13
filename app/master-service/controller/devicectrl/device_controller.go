package devicectrl

import (
	"api/app/master-service/domain/devicedomain"
	"api/app/master-service/model"
	"api/lib"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

type DeviceController struct {
	domain devicedomain.DeviceDomain
}

func New(domain devicedomain.DeviceDomain) *DeviceController {
	return &DeviceController{domain: domain}
}

// GetAll returns all devices
// @Summary Get all devices
// @Description Get all devices
// @Tags Device
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response{data=[]model.Device} "List of Devices"
// @Router /devices [get]
// @Security TokenKey
func (c *DeviceController) GetAll(ctx *fiber.Ctx) error {
	devices, err := c.domain.GetAll(ctx.Context())
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, devices)
}

// GetByID returns a device by ID
// @Summary Get a device by ID
// @Description Get a device by ID
// @Tags Device
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} resp.Response{data=model.Device} "Device data"
// @Router /devices/{id} [get]
// @Security TokenKey
func (c *DeviceController) GetByID(ctx *fiber.Ctx) error {
	device, err := c.domain.GetByID(ctx.Context(), lib.ParamsUUID(ctx))
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, device)
}

// Create creates a new device
// @Summary Create a new device
// @Description Create a new device
// @Tags Device
// @Accept json
// @Produce json
// @Param device body model.DeviceAPI true "Body Request"
// @Success 201 {object} resp.Response{data=model.Device} "Device data"
// @Router /devices [post]
// @Security TokenKey
func (c *DeviceController) Create(ctx *fiber.Ctx) error {
	api := new(model.DeviceAPI)
	if err := ctx.BodyParser(api); err != nil {
		return resp.ErrorHandler(ctx, err)
	}

	device, err := c.domain.Create(ctx.Context(), api)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}

	return resp.Created(ctx, device)
}

// Update updates a device
// @Summary Update a device
// @Description Update a device
// @Tags Device
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param device body model.DeviceUpdateRequest true "Body Request"
// @Success 200 {object} resp.Response{data=model.Device} "Device data"
// @Router /devices/{id} [put]
// @Security TokenKey
func (c *DeviceController) Update(ctx *fiber.Ctx) error {
	req := new(model.DeviceUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	device, err := c.domain.Update(ctx.Context(), lib.ParamsUUID(ctx), req)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, device)
}

// Delete deletes a device
// @Summary Delete a device
// @Description Delete a device
// @Tags Device
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} resp.Response "Success"
// @Router /devices/{id} [delete]
// @Security TokenKey
func (c *DeviceController) Delete(ctx *fiber.Ctx) error {
	if err := c.domain.Delete(ctx.Context(), lib.ParamsUUID(ctx)); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx)
}
