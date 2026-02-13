package devicerepo

import (
	"api/app/master-service/model"
	"api/entity"
	"context"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	entity.BaseRepository[model.Device]
	GetByCode(ctx context.Context, db *gorm.DB, code string) (*model.Device, error)
}

type deviceRepository struct {
	entity.BaseRepository[model.Device]
}

func New() DeviceRepository {
	return &deviceRepository{
		BaseRepository: entity.NewBaseRepository[model.Device](entity.Entity{
			Name: "device",
		}),
	}
}

func (r *deviceRepository) GetByCode(ctx context.Context, db *gorm.DB, code string) (*model.Device, error) {
	var device model.Device
	if err := db.WithContext(ctx).Where("code = ?", code).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
