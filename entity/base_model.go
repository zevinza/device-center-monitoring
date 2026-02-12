package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID         *uuid.UUID     `json:"id,omitempty" gorm:"primaryKey;unique;type:varchar(36);not null" format:"uuid" swaggerignore:"true"`
	CreatedAt  *time.Time     `json:"created_at,omitempty" gorm:"type:timestamptz" format:"date-time" swaggerignore:"true"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty" gorm:"type:timestamptz" format:"date-time" swaggerignore:"true"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
	Additional *string        `json:"additional,omitempty"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == nil {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		b.ID = &id
	}
	now := time.Now()
	if b.CreatedAt == nil {
		b.CreatedAt = &now
	}
	if b.UpdatedAt == nil {
		b.UpdatedAt = &now
	}
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	b.UpdatedAt = &now
	return nil
}
