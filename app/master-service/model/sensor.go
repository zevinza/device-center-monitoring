package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sensor is a managed sensor attached to a device (one device -> many sensors).
type Sensor struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DeviceID  uuid.UUID     `gorm:"type:uuid;not null;index" json:"device_id"`
	Device    Device        `gorm:"foreignKey:DeviceID;constraint:OnDelete:CASCADE" json:"-"`
	Name      string        `gorm:"type:varchar(255);not null" json:"name"`
	Unit      string        `gorm:"type:varchar(50)" json:"unit"`
	Type      string        `gorm:"type:varchar(100)" json:"type"`
	CreatedAt time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time     `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SensorCreateRequest struct {
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Type     string `json:"type"`
}

type SensorUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Unit *string `json:"unit,omitempty"`
	Type *string `json:"type,omitempty"`
}
