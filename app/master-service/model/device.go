package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device is a managed IoT device.
type Device struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DeviceCode  string          `gorm:"type:varchar(255);uniqueIndex;not null" json:"device_code"` // Public identifier for simulators
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description,omitempty"`
	Location    string          `gorm:"type:varchar(255)" json:"location,omitempty"`
	IsActive    bool            `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type DeviceCreateRequest struct {
	DeviceCode  *string `json:"device_code,omitempty"` // Optional, auto-generated if not provided
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
}

type DeviceUpdateRequest struct {
	DeviceCode  *string `json:"device_code,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type DeviceResponse struct {
	Device
	Sensors []Sensor `json:"sensors"`
}
