package sensordomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorrepo"
	"api/lib"
	"api/utils/resp"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type SensorDomain interface {
	GetAll(ctx context.Context, deviceID string) ([]model.Sensor, error)
	Create(ctx context.Context, req *model.SensorCreateRequest) (*model.Sensor, error)
	Update(ctx context.Context, id string, req *model.SensorUpdateRequest) (*model.Sensor, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Sensor, error)
}

type sensorDomain struct {
	sensorRepository sensorrepo.SensorRepository
	deviceRepository devicerepo.DeviceRepository
}

func New(sensorRepository sensorrepo.SensorRepository, deviceRepository devicerepo.DeviceRepository) SensorDomain {
	return &sensorDomain{
		sensorRepository: sensorRepository,
		deviceRepository: deviceRepository,
	}
}

func (d *sensorDomain) GetAll(ctx context.Context, deviceID string) ([]model.Sensor, error) {
	oid, err := lib.HexToObjectID(deviceID)
	if err != nil {
		return nil, resp.ErrorBadRequest(err.Error())
	}
	if _, err := d.deviceRepository.GetByID(ctx, oid); err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	sensors, err := d.sensorRepository.GetByDeviceID(ctx, oid, 100)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	return sensors, nil
}

func (d *sensorDomain) GetByID(ctx context.Context, id string) (*model.Sensor, error) {
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	sensor, err := d.sensorRepository.GetByID(ctx, oid)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return sensor, nil
}

func (d *sensorDomain) Create(ctx context.Context, req *model.SensorCreateRequest) (*model.Sensor, error) {
	if req.DeviceID.IsZero() || strings.TrimSpace(req.Name) == "" {
		return nil, resp.ErrorBadRequest("invalid request")
	}
	if _, err := d.deviceRepository.GetByID(ctx, req.DeviceID); err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	s := &model.Sensor{
		DeviceID: req.DeviceID,
		Name:     strings.TrimSpace(req.Name),
		Unit:     strings.TrimSpace(req.Unit),
		Type:     strings.TrimSpace(req.Type),
	}
	if err := d.sensorRepository.Create(ctx, s); err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return s, nil
}

func (d *sensorDomain) Update(ctx context.Context, id string, req *model.SensorUpdateRequest) (*model.Sensor, error) {
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	update := bson.M{}
	if req.Name != nil {
		update["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Unit != nil {
		update["unit"] = strings.TrimSpace(*req.Unit)
	}
	if req.Type != nil {
		update["type"] = strings.TrimSpace(*req.Type)
	}

	sensor, err := d.sensorRepository.Update(ctx, oid, update)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}

	return sensor, nil
}

func (d *sensorDomain) Delete(ctx context.Context, id string) error {
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if _, err := d.sensorRepository.GetByID(ctx, oid); err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if err := d.sensorRepository.Delete(ctx, oid); err != nil {
		return resp.ErrorInternal(err.Error())
	}
	return nil
}
