package sensordomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorrepo"
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidSensor = errors.New("invalid sensor payload")

type SensorDomain interface {
	Create(ctx context.Context, req *model.SensorCreateRequest) (*model.Sensor, error)
	Update(ctx context.Context, id string, req *model.SensorUpdateRequest) (*model.Sensor, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Sensor, error)
	ListByDeviceID(ctx context.Context, deviceID string, limit int64) ([]model.Sensor, error)
}

type sensorDomain struct {
	sensorRepo sensorrepo.SensorRepository
	deviceRepo devicerepo.DeviceRepository
}

func New(sensorRepo sensorrepo.SensorRepository, deviceRepo devicerepo.DeviceRepository) SensorDomain {
	return &sensorDomain{sensorRepo: sensorRepo, deviceRepo: deviceRepo}
}

func (d *sensorDomain) Create(ctx context.Context, req *model.SensorCreateRequest) (*model.Sensor, error) {
	if req.DeviceID.IsZero() || strings.TrimSpace(req.Name) == "" {
		return nil, ErrInvalidSensor
	}
	// Ensure device exists (convert ObjectID to hex string for lookup)
	deviceIDHex := req.DeviceID.Hex()
	if _, err := d.deviceRepo.FindByID(ctx, deviceIDHex); err != nil {
		return nil, err
	}
	s := &model.Sensor{
		DeviceID: req.DeviceID,
		Name:     strings.TrimSpace(req.Name),
		Unit:     strings.TrimSpace(req.Unit),
		Type:     strings.TrimSpace(req.Type),
	}
	if err := d.sensorRepo.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (d *sensorDomain) Update(ctx context.Context, id string, req *model.SensorUpdateRequest) (*model.Sensor, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return nil, err
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
	return d.sensorRepo.Update(ctx, oid, update)
}

func (d *sensorDomain) Delete(ctx context.Context, id string) error {
	oid, err := parseObjectID(id)
	if err != nil {
		return err
	}
	return d.sensorRepo.Delete(ctx, oid)
}

func (d *sensorDomain) GetByID(ctx context.Context, id string) (*model.Sensor, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return nil, err
	}
	return d.sensorRepo.FindByID(ctx, oid)
}

func (d *sensorDomain) ListByDeviceID(ctx context.Context, deviceID string, limit int64) ([]model.Sensor, error) {
	return d.sensorRepo.FindByDeviceID(ctx, deviceID, limit)
}

func parseObjectID(hex string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(strings.TrimSpace(hex))
}
