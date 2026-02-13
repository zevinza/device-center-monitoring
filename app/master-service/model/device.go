package model

import "api/entity"

// Device is a managed IoT device.
type Device struct {
	entity.Base
	DeviceAPI
}
type DeviceAPI struct {
	Code        string  `json:"code" gorm:"type:varchar(255);uniqueIndex;not null"` // Public identifier for simulators
	Name        string  `json:"name" gorm:"type:varchar(255);not null"`
	Description *string `json:"description,omitempty" gorm:"type:text"`
	Location    *string `json:"location,omitempty" gorm:"type:varchar(255)"`
	IsActive    bool    `json:"is_active" gorm:"default:true;not null"`
}

type DeviceUpdateRequest struct {
	Code        *string `json:"code,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Location    *string `json:"location,omitempty"`
}

type DeviceResponse struct {
	Device
	Sensors []Sensor `json:"sensors"`
}
