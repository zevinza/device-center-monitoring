package sensorcategorydomain

import (
	"api/app/master-service/model"
	"api/app/master-service/repository/sensorcategory"
	"api/lib"
	"api/utils/resp"
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SensorCategoryDomain interface {
	GetAll(ctx context.Context) ([]model.SensorCategoryResponse, error)
	GetByID(ctx context.Context, id *uuid.UUID) (*model.SensorCategoryResponse, error)
	Create(ctx context.Context, api *model.SensorCategoryRequest) (*model.SensorCategoryResponse, error)
	Update(ctx context.Context, id *uuid.UUID, req *model.SensorCategoryRequest) (*model.SensorCategoryResponse, error)
	Delete(ctx context.Context, id *uuid.UUID) error
}

type sensorCategoryDomain struct {
	db                       *gorm.DB
	sensorCategoryRepository sensorcategory.SensorCategoryRepository
}

func New(db *gorm.DB, sensorCategoryRepository sensorcategory.SensorCategoryRepository) SensorCategoryDomain {
	return &sensorCategoryDomain{
		db:                       db,
		sensorCategoryRepository: sensorCategoryRepository,
	}
}

func (d *sensorCategoryDomain) GetAll(ctx context.Context) ([]model.SensorCategoryResponse, error) {
	categories, err := d.sensorCategoryRepository.GetAll(ctx, d.db)
	if err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	var categoriesResponse []model.SensorCategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, lib.Rev(category.ToSensorCategoryResponse()))
	}
	return categoriesResponse, nil
}

func (d *sensorCategoryDomain) GetByID(ctx context.Context, id *uuid.UUID) (*model.SensorCategoryResponse, error) {
	category, err := d.sensorCategoryRepository.GetByID(ctx, d.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor category not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}
	return category.ToSensorCategoryResponse(), nil
}

func (d *sensorCategoryDomain) Create(ctx context.Context, api *model.SensorCategoryRequest) (*model.SensorCategoryResponse, error) {
	if len(api.Units) == 0 {
		return nil, resp.ErrorBadRequest("units are required")
	}

	if !slices.Contains(api.Units, lib.Rev(api.DefaultUnit)) {
		return nil, resp.ErrorBadRequest("default unit is not in units")
	}

	category := model.SensorCategory{}
	lib.Merge(api, &category)

	if err := d.sensorCategoryRepository.Create(ctx, d.db, &category); err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return category.ToSensorCategoryResponse(), nil
}

func (d *sensorCategoryDomain) Update(ctx context.Context, id *uuid.UUID, req *model.SensorCategoryRequest) (*model.SensorCategoryResponse, error) {
	category, err := d.sensorCategoryRepository.GetByID(ctx, d.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrorNotFound("Sensor category not found")
		}
		return nil, resp.ErrorInternal(err.Error())
	}

	lib.Merge(req, category)
	if !slices.Contains(req.Units, lib.Rev(req.DefaultUnit)) {
		return nil, resp.ErrorBadRequest("default unit is not in units")
	}
	category.Units = strings.Join(req.Units, "|")
	if err := d.sensorCategoryRepository.Update(ctx, d.db, category, id); err != nil {
		return nil, resp.ErrorInternal(err.Error())
	}
	return category.ToSensorCategoryResponse(), nil
}

func (d *sensorCategoryDomain) Delete(ctx context.Context, id *uuid.UUID) error {
	if _, err := d.sensorCategoryRepository.GetByID(ctx, d.db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.ErrorNotFound("Sensor category not found")
		}
		return resp.ErrorInternal(err.Error())
	}
	if err := d.sensorCategoryRepository.Delete(ctx, d.db, id); err != nil {
		return resp.ErrorInternal(err.Error())
	}
	return nil
}
