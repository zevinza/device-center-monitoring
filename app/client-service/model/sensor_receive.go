package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorReceiveRequest struct {
	SensorReading
	Device *Device `json:"device"`
	Sensor *Sensor `json:"sensor,omitempty"`
}

type SensorReading struct {
	ID         primitive.ObjectID `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	SensorID   *uuid.UUID         `json:"sensor_id"`
	Value      any                `json:"value"`
	Timestamp  *time.Time         `json:"timestamp"`
	Status     string             `json:"status"`
	RetryCount int                `json:"retry_count"`
}

type Device struct {
	ID          *uuid.UUID `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Additional  *string    `json:"additional,omitempty"`
}

type Sensor struct {
	ID         *uuid.UUID `json:"id"`
	DeviceID   *uuid.UUID `json:"device_id"`
	CategoryID *uuid.UUID `json:"category_id"`
	Name       string     `json:"name"`
	Unit       string     `json:"unit"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	Additional *string    `json:"additional,omitempty"`
}
