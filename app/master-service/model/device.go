package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Device is a managed IoT device.
type Device struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	DeviceCode  string             `bson:"device_code" json:"device_code"` // Public identifier for simulators
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
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
