package model

import (
	"api/entity"

	"github.com/google/uuid"
)

// Sensor is a managed sensor attached to a device (one device -> many sensors).
type Sensor struct {
	entity.Base
	SensorAPI
	Device   *Device         `json:"device" gorm:"foreignKey:DeviceID;constraint:OnDelete:CASCADE"`
	Category *SensorCategory `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Readings []SensorReading `json:"readings" gorm:"-"`
}
type SensorAPI struct {
	DeviceID   *uuid.UUID `json:"device_id" gorm:"type:uuid;not null;index"`
	CategoryID *uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	Name       string     `json:"name" gorm:"type:varchar(255);not null"`
	Unit       string     `json:"unit" gorm:"type:varchar(50)"`
}

type SensorUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Unit *string `json:"unit,omitempty"`
	Type *string `json:"type,omitempty"`
}
