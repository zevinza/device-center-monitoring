package devicedomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorrepo"
	"api/lib"
	"api/utils/resp"
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceDomain interface {
	GetAll(ctx context.Context) ([]model.Device, error)
	GetByID(ctx context.Context, id *uuid.UUID) (*model.Device, error)
	Create(ctx context.Context, req *model.DeviceAPI) (*model.Device, error)
	Update(ctx context.Context, id *uuid.UUID, req *model.DeviceUpdateRequest) (*model.Device, error)
	Delete(ctx context.Context, id *uuid.UUID) error
}

type deviceDomain struct {
	db               *gorm.DB
	deviceRepository devicerepo.DeviceRepository
	sensorRepository sensorrepo.SensorRepository
}

func New(db *gorm.DB, deviceRepository devicerepo.DeviceRepository, sensorRepository sensorrepo.SensorRepository) DeviceDomain {
	return &deviceDomain{
		db:               db,
		deviceRepository: deviceRepository,
		sensorRepository: sensorRepository,
	}
}

func (d *deviceDomain) GetAll(ctx context.Context) ([]model.Device, error) {
	devices, err := d.deviceRepository.GetAll(ctx, d.db)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return devices, nil
}

func (d *deviceDomain) GetByID(ctx context.Context, id *uuid.UUID) (*model.Device, error) {
	device, err := d.deviceRepository.GetByID(ctx, d.db, id)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}

	return device, nil
}

func (d *deviceDomain) Create(ctx context.Context, req *model.DeviceAPI) (*model.Device, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, resp.ErrorBadRequest("device name is required")
	}

	device := model.Device{}
	lib.Merge(req, &device)

	if err := d.deviceRepository.Create(ctx, d.db, &device); err != nil {
		if err == gorm.ErrDuplicatedKey {
			return nil, resp.ErrorConflict("device code already exists")
		}
		return nil, resp.ErrorInternal(err.Error())
	}
	return &device, nil
}

func (d *deviceDomain) Update(ctx context.Context, id *uuid.UUID, req *model.DeviceUpdateRequest) (*model.Device, error) {
	device, err := d.deviceRepository.GetByID(ctx, d.db, id)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	lib.Merge(req, &device)
	if err := d.deviceRepository.Update(ctx, d.db, device, id); err != nil {
		if err == gorm.ErrDuplicatedKey {
			return nil, resp.ErrorConflict("device code already exists")
		}
		return nil, resp.ErrorInternal(err.Error())
	}
	return device, nil
}

func (d *deviceDomain) Delete(ctx context.Context, id *uuid.UUID) error {
	if _, err := d.deviceRepository.GetByID(ctx, d.db, id); err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if err := d.deviceRepository.Delete(ctx, d.db, id); err != nil {
		return resp.ErrorInternal(err.Error())
	}
	return nil
}
