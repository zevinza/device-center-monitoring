package sensorrepo

import (
	"api/app/master-service/model"
	"api/entity"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SensorRepository interface {
	entity.BaseRepository[model.Sensor]
	GetByDeviceID(ctx context.Context, db *gorm.DB, deviceID *uuid.UUID) ([]model.Sensor, error)
	GetByNameAndDeviceID(ctx context.Context, db *gorm.DB, name string, deviceID *uuid.UUID, excludeID *uuid.UUID) (*model.Sensor, error)
}

type sensorRepository struct {
	entity.BaseRepository[model.Sensor]
}

func New() SensorRepository {
	return &sensorRepository{
		BaseRepository: entity.NewBaseRepository[model.Sensor](entity.Entity{
			Name: "sensor",
		}),
	}
}

func (r *sensorRepository) GetAll(ctx context.Context, db *gorm.DB) ([]model.Sensor, error) {
	var sensors []model.Sensor
	if err := db.WithContext(ctx).Joins("Device").Joins("Category").Find(&sensors).Error; err != nil {
		return nil, err
	}
	return sensors, nil
}

func (r *sensorRepository) GetByID(ctx context.Context, db *gorm.DB, id *uuid.UUID) (*model.Sensor, error) {
	var sensor model.Sensor
	if err := db.WithContext(ctx).Where("sensor.id = ?", id).Joins("Device").Joins("Category").First(&sensor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sensor not found")
		}
		return nil, err
	}
	return &sensor, nil
}

func (r *sensorRepository) GetByDeviceID(ctx context.Context, db *gorm.DB, deviceID *uuid.UUID) ([]model.Sensor, error) {
	var sensors []model.Sensor
	query := db.WithContext(ctx).Where("device_id = ?", deviceID).Joins("Device").Joins("Category").Find(&sensors)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sensors not found")
		}
		return nil, query.Error
	}
	return sensors, nil
}

func (r *sensorRepository) GetByNameAndDeviceID(ctx context.Context, db *gorm.DB, name string, deviceID *uuid.UUID, excludeID *uuid.UUID) (*model.Sensor, error) {
	var sensor model.Sensor
	query := db.WithContext(ctx).Where("name = ? AND device_id = ?", name, deviceID)
	if excludeID != nil {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.First(&sensor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sensor not found")
		}
		return nil, err
	}
	return &sensor, nil
}
