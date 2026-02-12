package devicectrl

import (
	"api/app/master-service/domain/devicedomain"
	"api/app/master-service/model"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

type DeviceController struct {
	domain devicedomain.DeviceDomain
}

func New(domain devicedomain.DeviceDomain) *DeviceController {
	return &DeviceController{domain: domain}
}

func (ctl *DeviceController) Create(c *fiber.Ctx) error {
	req := new(model.DeviceCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	device, err := ctl.domain.Create(c.Context(), req)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	return resp.Created(c, device)
}

func (ctl *DeviceController) GetAll(c *fiber.Ctx) error {
	devices, err := ctl.domain.GetAll(c.Context())
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorInternal(err.Error()))
	}
	return resp.OK(c, devices)
}

func (ctl *DeviceController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	device, err := ctl.domain.GetByID(c.Context(), id)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorNotFound(err.Error()))
	}
	return resp.OK(c, device)
}

func (ctl *DeviceController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(model.DeviceUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	device, err := ctl.domain.Update(c.Context(), id, req)
	if err != nil {
		return resp.ErrorHandler(c, resp.ErrorBadRequest(err.Error()))
	}
	return resp.OK(c, device)
}

func (ctl *DeviceController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := ctl.domain.Delete(c.Context(), id); err != nil {
		return resp.ErrorHandler(c, resp.ErrorNotFound(err.Error()))
	}
	return resp.OK(c)
}
