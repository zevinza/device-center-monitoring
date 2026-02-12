package sensorctrl

import (
	"api/app/master-service/domain/sensorreadingdomain"
	"api/app/master-service/model"
	"api/utils/resp"

	"github.com/gofiber/fiber/v2"
)

type SensorIngestController struct {
	domain sensorreadingdomain.SensorReadingDomain
}

func NewSensorIngestController(
	domain sensorreadingdomain.SensorReadingDomain,
) *SensorIngestController {
	return &SensorIngestController{
		domain: domain,
	}
}

// Ingest ingests a sensor reading
// @Summary Ingest a sensor reading
// @Description Ingest a sensor reading
// @Tags Sensor Ingest
// @Accept json
// @Produce json
// @Param sensor_reading body model.SensorIngestRequest true "Sensor reading data"
// @Success 200 {object} resp.Response "Success"
// @Router /sensors/ingest [post]
// @Security TokenKey
func (c *SensorIngestController) Ingest(ctx *fiber.Ctx) error {
	api := new(model.SensorIngestRequest)
	if err := ctx.BodyParser(api); err != nil {
		return resp.ErrorHandler(ctx, resp.ErrorBadRequest(err.Error()))
	}

	if err := c.domain.Ingest(ctx.Context(), api); err != nil {
		return resp.ErrorHandler(ctx, err)
	}

	return resp.OK(ctx)
}
