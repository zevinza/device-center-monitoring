package devicedomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/lib"
	"api/utils/resp"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceDomain interface {
	GetAll(ctx context.Context) ([]model.Device, error)
	GetByID(ctx context.Context, id string) (*model.Device, error)
	Create(ctx context.Context, req *model.DeviceCreateRequest) (*model.Device, error)
	Update(ctx context.Context, id string, req *model.DeviceUpdateRequest) (*model.Device, error)
	Delete(ctx context.Context, id string) error
}

type deviceDomain struct {
	deviceRepository devicerepo.DeviceRepository
}

func New(deviceRepository devicerepo.DeviceRepository) DeviceDomain {
	return &deviceDomain{deviceRepository: deviceRepository}
}

func (d *deviceDomain) GetAll(ctx context.Context) ([]model.Device, error) {
	devices, err := d.deviceRepository.FindAll(ctx, 0)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return devices, nil
}

func (d *deviceDomain) GetByID(ctx context.Context, id string) (*model.Device, error) {
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	device, err := d.deviceRepository.GetByID(ctx, oid)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return device, nil
}

func (d *deviceDomain) Create(ctx context.Context, req *model.DeviceCreateRequest) (*model.Device, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("device name is required")
	}
	now := time.Now().UTC()

	deviceCode := ""
	if req.DeviceCode != nil && strings.TrimSpace(*req.DeviceCode) != "" {
		deviceCode = strings.TrimSpace(*req.DeviceCode)
		// Check if device_code already exists
		if existing, _ := d.deviceRepository.FindByCode(ctx, deviceCode); existing != nil {
			return nil, errors.New("device_code already exists")
		}
	} else {
		deviceCode = "DEV-" + primitive.NewObjectID().Hex()[12:]
	}

	device := &model.Device{
		ID:          primitive.NewObjectID(),
		DeviceCode:  deviceCode,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := d.deviceRepository.Create(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (d *deviceDomain) Update(ctx context.Context, id string, req *model.DeviceUpdateRequest) (*model.Device, error) {
	update := map[string]any{}
	if req.DeviceCode != nil {
		deviceCode := strings.TrimSpace(*req.DeviceCode)
		if deviceCode != "" {
			// Check if device_code already exists (excluding current device)
			if existing, _ := d.deviceRepository.FindByCode(ctx, deviceCode); existing != nil && existing.ID.Hex() != id {
				return nil, errors.New("device_code already exists")
			}
			update["device_code"] = deviceCode
		}
	}
	if req.Name != nil {
		update["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		update["description"] = strings.TrimSpace(*req.Description)
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	device, err := d.deviceRepository.Update(ctx, oid, update)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return device, nil
}

func (d *deviceDomain) Delete(ctx context.Context, id string) error {
	oid, err := lib.HexToObjectID(id)
	if err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if _, err := d.deviceRepository.GetByID(ctx, oid); err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if err := d.deviceRepository.Delete(ctx, oid); err != nil {
		return resp.ErrorInternal(err.Error())
	}
	return nil
}
