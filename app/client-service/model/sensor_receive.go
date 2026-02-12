package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorReceiveRequest struct {
	SensorReading
	Device *Device `json:"device"`
	Sensor *Sensor `json:"sensor,omitempty"`
}

type SensorReading struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	DeviceID   primitive.ObjectID `bson:"device_id" json:"device_id"`
	SensorID   primitive.ObjectID `bson:"sensor_id" json:"sensor_id"`
	Value      string             `bson:"value" json:"value"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Status     string             `bson:"status" json:"status"`
	RetryCount int                `bson:"retry_count" json:"retry_count"`
}

type Device struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	DeviceCode  string             `bson:"device_code" json:"device_code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Sensor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID  primitive.ObjectID `bson:"device_id" json:"device_id"`
	Name      string             `bson:"name" json:"name"`
	Unit      string             `bson:"unit" json:"unit"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
