package model

import "time"

// SensorDataRequest is the payload received from master service
type SensorDataRequest struct {
	ReadingID string      `json:"reading_id"`
	DeviceID  string      `json:"device_id"`
	Device    *DeviceInfo `json:"device"`
	Sensor    *SensorInfo `json:"sensor,omitempty"`
	Unit      string      `json:"unit"`
	Value     float64     `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// DeviceInfo contains device details
type DeviceInfo struct {
	ID          string `json:"id"`          // ObjectID hex string
	DeviceCode  string `json:"device_code"` // Public identifier
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// SensorInfo contains sensor details (if available)
type SensorInfo struct {
	ID   string `json:"id"` // ObjectID hex string
	Name string `json:"name"`
	Unit string `json:"unit"`
	Type string `json:"type,omitempty"`
}
