package model

type QueueMessage struct {
	ReadingID string `json:"reading_id"`
}

type SensorIngestResult struct {
	SensorReading
	Device *Device `json:"device"`
	Sensor *Sensor `json:"sensor,omitempty"`
}
