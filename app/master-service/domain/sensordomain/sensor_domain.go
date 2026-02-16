package sensordomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/devicerepo"
	"api/app/master-service/repository/sensorcategory"
	"api/app/master-service/repository/sensorreadingrepo"
	"api/app/master-service/repository/sensorrepo"
	"api/lib"
	"api/utils/resp"
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type SensorDomain interface {
	GetAll(ctx context.Context, deviceID *uuid.UUID) ([]model.Sensor, error)
	GetByID(ctx context.Context, limit int64, deviceID, id *uuid.UUID) (*model.Sensor, error)
	Create(ctx context.Context, api *model.SensorAPI) (*model.Sensor, error)
	Update(ctx context.Context, deviceID, id *uuid.UUID, api *model.SensorUpdateRequest) (*model.Sensor, error)
	Delete(ctx context.Context, deviceID, id *uuid.UUID) error
}

type sensorDomain struct {
	db *gorm.DB

	sensorRepository         sensorrepo.SensorRepository
	deviceRepository         devicerepo.DeviceRepository
	sensorReadingRepository  sensorreadingrepo.SensorReadingRepository
	sensorCategoryRepository sensorcategory.SensorCategoryRepository
}

func New(db *gorm.DB, sensorRepository sensorrepo.SensorRepository, deviceRepository devicerepo.DeviceRepository, sensorReadingRepository sensorreadingrepo.SensorReadingRepository, sensorCategoryRepository sensorcategory.SensorCategoryRepository) SensorDomain {
	return &sensorDomain{
		db:                       db,
		sensorRepository:         sensorRepository,
		deviceRepository:         deviceRepository,
		sensorReadingRepository:  sensorReadingRepository,
		sensorCategoryRepository: sensorCategoryRepository,
	}
}

func (d *sensorDomain) GetAll(ctx context.Context, deviceID *uuid.UUID) ([]model.Sensor, error) {
	if _, err := d.deviceRepository.GetByID(ctx, d.db, deviceID); err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	sensors, err := d.sensorRepository.GetByDeviceID(ctx, d.db, deviceID)
	if err != nil {
		return nil, resp.ErrorNotFound(err.Error())
	}
	return sensors, nil
}

func (d *sensorDomain) GetByID(ctx context.Context, limit int64, deviceID, id *uuid.UUID) (*model.Sensor, error) {
	if _, err := d.deviceRepository.GetByID(ctx, d.db, deviceID); err != nil {
		return nil, resp.ErrorNotFound("Device not found")
	}
	sensor, err := d.sensorRepository.GetByID(ctx, d.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}

	readings, err := d.sensorReadingRepository.GetBySensorID(ctx, lib.StringUUID(sensor.ID), limit)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, resp.ErrorInternal(err.Error())
	}
	sensor.Readings = readings
	return sensor, nil
}

func (d *sensorDomain) Create(ctx context.Context, api *model.SensorAPI) (*model.Sensor, error) {
	if _, err := d.deviceRepository.GetByID(ctx, d.db, api.DeviceID); err != nil {
		return nil, resp.ErrorNotFound("Device not found")
	}

	// Validate category exists
	category, err := d.sensorCategoryRepository.GetByID(ctx, d.db, api.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor category not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}

	sensor := model.Sensor{}
	lib.Merge(api, &sensor)
	if sensor.Name == "" {
		return nil, resp.ErrorBadRequest("sensor name is required")
	}

	// Validate that sensor name is unique within the device
	existingSensor, err := d.sensorRepository.GetByNameAndDeviceID(ctx, d.db, sensor.Name, api.DeviceID, nil)
	if err == nil && existingSensor != nil {
		return nil, resp.ErrorConflict("sensor name already exists for this device")
	}
	// If error is not "sensor not found", it's a database error
	if err != nil && err.Error() != "sensor not found" {
		return nil, resp.ErrorInternal(err.Error())
	}

	// Validate and set unit
	unitNames := category.GetUnitNames()
	if sensor.Unit == "" {
		sensor.Unit = category.DefaultUnit
	} else {
		if !slices.Contains(unitNames, sensor.Unit) {
			return nil, resp.ErrorBadRequest("unit must be one of the available units for this category")
		}
	}

	if err := d.sensorRepository.Create(ctx, d.db, &sensor); err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return &sensor, nil
}

func (d *sensorDomain) Update(ctx context.Context, deviceID, id *uuid.UUID, api *model.SensorUpdateRequest) (*model.Sensor, error) {
	if _, err := d.deviceRepository.GetByID(ctx, d.db, deviceID); err != nil {
		return nil, resp.ErrorNotFound("Device not found")
	}

	// Get existing sensor to check category
	existingSensor, err := d.sensorRepository.GetByID(ctx, d.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}

	// Get category for validation
	categoryID := existingSensor.CategoryID
	category, err := d.sensorCategoryRepository.GetByID(ctx, d.db, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor category not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}

	sensor := model.Sensor{}
	lib.Merge(api, &sensor)

	// Validate that sensor name is unique within the device (if name is being updated)
	if api.Name != nil && *api.Name != "" {
		existingSensor, err := d.sensorRepository.GetByNameAndDeviceID(ctx, d.db, *api.Name, deviceID, id)
		if err == nil && existingSensor != nil {
			return nil, resp.ErrorConflict("sensor name already exists for this device")
		}
		// If error is not "sensor not found", it's a database error
		if err != nil && err.Error() != "sensor not found" {
			return nil, resp.ErrorInternal(err.Error())
		}
	}

	if api.Unit != nil {
		unitNames := category.GetUnitNames()
		if *api.Unit == "" {
			sensor.Unit = category.DefaultUnit
		} else {
			if !slices.Contains(unitNames, *api.Unit) {
				return nil, resp.ErrorBadRequest("unit must be one of the available units for this category")
			}
			sensor.Unit = *api.Unit
		}
	}

	if err := d.sensorRepository.Update(ctx, d.db, &sensor, id); err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return &sensor, nil
}

func (d *sensorDomain) Delete(ctx context.Context, deviceID, id *uuid.UUID) error {
	if _, err := d.sensorRepository.GetByID(ctx, d.db, id); err != nil {
		return resp.ErrorNotFound(err.Error())
	}
	if err := d.sensorRepository.Delete(ctx, d.db, id); err != nil {
		return resp.ErrorInternal(err.Error())
	}
	return nil
}
