package entity

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	GetAll(ctx context.Context, db *gorm.DB) ([]T, error)
	GetByID(ctx context.Context, db *gorm.DB, id *uuid.UUID) (*T, error)
	GetByColumn(ctx context.Context, db *gorm.DB, column string, value any) (*T, error)
	Create(ctx context.Context, tx *gorm.DB, data *T) error
	Update(ctx context.Context, tx *gorm.DB, data *T, id *uuid.UUID) error
	Save(ctx context.Context, tx *gorm.DB, values map[string]any, id *uuid.UUID) error
	Delete(ctx context.Context, tx *gorm.DB, id *uuid.UUID) error
}

var _ BaseRepository[any] = &baseRepository[any]{}

type baseRepository[T any] struct {
	entity Entity
}

// Entity describes a database entity (table or view)
type Entity struct {
	Name   string
	IsView bool
}

// NewBaseRepository creates a new BaseRepository instance with the given database connection
func NewBaseRepository[T any](entity Entity) BaseRepository[T] {
	return &baseRepository[T]{
		entity: entity,
	}
}

func (r *baseRepository[T]) GetAll(ctx context.Context, db *gorm.DB) ([]T, error) {
	var data []T
	if err := db.WithContext(ctx).Table(r.entity.Name).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *baseRepository[T]) GetByID(ctx context.Context, db *gorm.DB, id *uuid.UUID) (*T, error) {
	var data T
	if err := db.WithContext(ctx).Table(r.entity.Name).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *baseRepository[T]) GetByColumn(ctx context.Context, db *gorm.DB, column string, value any) (*T, error) {
	var data T
	if !IsValidColumn[T](db, column) {
		return nil, errors.New("column is not valid")
	}
	if err := db.WithContext(ctx).Table(r.entity.Name).Where(column+" = ?", value).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *baseRepository[T]) Create(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Table(r.entity.Name).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, tx *gorm.DB, data *T, id *uuid.UUID) error {
	if err := tx.WithContext(ctx).Table(r.entity.Name).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *baseRepository[T]) Save(ctx context.Context, tx *gorm.DB, values map[string]any, id *uuid.UUID) error {
	if err := tx.WithContext(ctx).Table(r.entity.Name).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (r *baseRepository[T]) Delete(ctx context.Context, tx *gorm.DB, id *uuid.UUID) error {
	t := new(T)
	return tx.WithContext(ctx).Table(r.entity.Name).Delete(t, id).Error
}

func IsValidColumn[T any](db *gorm.DB, column string) bool {
	var model T
	return db.Migrator().HasColumn(&model, column)
}
