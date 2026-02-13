package sensorcategory

import (
	"api/app/master-service/model"
	"api/entity"
)

type SensorCategoryRepository interface {
	entity.BaseRepository[model.SensorCategory]
}

type sensorCategoryRepository struct {
	entity.BaseRepository[model.SensorCategory]
}

func New() SensorCategoryRepository {
	return &sensorCategoryRepository{
		BaseRepository: entity.NewBaseRepository[model.SensorCategory](entity.Entity{
			Name: "sensor_category",
		}),
	}
}
