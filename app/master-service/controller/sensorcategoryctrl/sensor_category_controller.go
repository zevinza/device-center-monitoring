package sensorcategoryctrl

import (
	"api/app/master-service/domain/sensorcategorydomain"
	"api/app/master-service/model"
	"api/lib"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

type SensorCategoryController struct {
	sensorCategoryDomain sensorcategorydomain.SensorCategoryDomain
}

func New(sensorCategoryDomain sensorcategorydomain.SensorCategoryDomain) *SensorCategoryController {
	return &SensorCategoryController{sensorCategoryDomain: sensorCategoryDomain}
}

// GetAll gets all sensor categories
// @Summary Get all sensor categories
// @Description Get all sensor categories
// @Tags SensorCategory
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response{data=[]model.SensorCategory} "List of Sensor Categories"
// @Router /sensor-categories [get]
// @Security TokenKey
func (c *SensorCategoryController) GetAll(ctx *fiber.Ctx) error {
	categories, err := c.sensorCategoryDomain.GetAll(ctx.Context())
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, categories)
}

// GetByID gets a sensor category by id
// @Summary Get a sensor category by id
// @Description Get a sensor category by id
// @Tags SensorCategory
// @Accept json
// @Produce json
// @Param id path string true "Sensor Category ID"
// @Success 200 {object} resp.Response{data=model.SensorCategory} "Sensor Category data"
// @Router /sensor-categories/{id} [get]
// @Security TokenKey
func (c *SensorCategoryController) GetByID(ctx *fiber.Ctx) error {
	category, err := c.sensorCategoryDomain.GetByID(ctx.Context(), lib.ParamsUUID(ctx))
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, category)
}

// Create creates a new sensor category
// @Summary Create a new sensor category
// @Description Create a new sensor category
// @Tags SensorCategory
// @Accept json
// @Produce json
// @Param sensor_category body model.SensorCategoryRequest true "Body Request"
// @Success 201 {object} resp.Response{data=model.SensorCategory} "Sensor Category data"
// @Router /sensor-categories [post]
// @Security TokenKey
func (c *SensorCategoryController) Create(ctx *fiber.Ctx) error {
	api := new(model.SensorCategoryRequest)
	if err := ctx.BodyParser(api); err != nil {
		return resp.ErrorHandler(ctx, err)
	}

	category, err := c.sensorCategoryDomain.Create(ctx.Context(), api)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.Created(ctx, category)
}

// Update updates a sensor category
// @Summary Update a sensor category
// @Description Update a sensor category
// @Tags SensorCategory
// @Accept json
// @Produce json
// @Param id path string true "Sensor Category ID"
// @Param sensor_category body model.SensorCategoryRequest true "Body Request"
// @Success 200 {object} resp.Response{data=model.SensorCategory} "Sensor Category data"
// @Router /sensor-categories/{id} [put]
// @Security TokenKey
func (c *SensorCategoryController) Update(ctx *fiber.Ctx) error {
	req := new(model.SensorCategoryRequest)
	if err := ctx.BodyParser(req); err != nil {
		return resp.ErrorHandler(ctx, err)
	}

	category, err := c.sensorCategoryDomain.Update(ctx.Context(), lib.ParamsUUID(ctx), req)
	if err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx, category)
}

// Delete deletes a sensor category
// @Summary Delete a sensor category
// @Description Delete a sensor category
// @Tags SensorCategory
// @Accept json
// @Produce json
// @Param id path string true "Sensor Category ID"
// @Success 200 {object} resp.Response "Success"
// @Router /sensor-categories/{id} [delete]
// @Security TokenKey
func (c *SensorCategoryController) Delete(ctx *fiber.Ctx) error {
	if err := c.sensorCategoryDomain.Delete(ctx.Context(), lib.ParamsUUID(ctx)); err != nil {
		return resp.ErrorHandler(ctx, err)
	}
	return resp.OK(ctx)
}
