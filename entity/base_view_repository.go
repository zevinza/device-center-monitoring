package entity

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseViewRepository[T any] interface {
	GetPaginated(ctx context.Context, db *gorm.DB, filter *Filter) ([]T, int64, error)
	GetFiltered(ctx context.Context, db *gorm.DB, filter *Filter) ([]T, int64, error)
	GetDetail(ctx context.Context, db *gorm.DB, id *uuid.UUID) (*T, error)
}

type baseViewRepository[T any] struct {
	entity Entity
}

func NewBaseViewRepository[T any](entity Entity) BaseViewRepository[T] {
	return &baseViewRepository[T]{
		entity: entity,
	}
}

func (r *baseViewRepository[T]) GetPaginated(ctx context.Context, db *gorm.DB, filter *Filter) ([]T, int64, error) {
	query := db.WithContext(ctx).Table(r.entity.Name)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if len(filter.Sort) > 0 {
		sorts, err := r.buildOrderClause(db, filter)
		if err != nil {
			return nil, 0, err
		}
		query = query.Order(sorts)
	}

	if filter.Page >= 0 {
		if filter.Page > 0 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
	}
	var data []T
	if err := query.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *baseViewRepository[T]) GetFiltered(ctx context.Context, db *gorm.DB, filter *Filter) ([]T, int64, error) {
	var data []T
	if err := db.WithContext(ctx).Table(r.entity.Name).Where(filter.Filter).Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, int64(len(data)), nil
}

func (r *baseViewRepository[T]) GetDetail(ctx context.Context, db *gorm.DB, id *uuid.UUID) (*T, error) {
	var data T
	if err := db.WithContext(ctx).Table(r.entity.Name).Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *baseViewRepository[T]) buildOrderClause(db *gorm.DB, filter *Filter) (clause.OrderBy, error) {
	var sorts []clause.OrderByColumn
	for _, sort := range filter.Sort {
		column := strings.TrimPrefix(sort, "-")
		if !IsValidColumn[T](db, column) {
			continue
		}
		sorts = append(sorts, clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   strings.HasPrefix(sort, "-"),
		})
	}
	return clause.OrderBy{Columns: sorts}, nil
}
