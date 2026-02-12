package config

var Environment = map[string]any{
	"ENVIRONMENT":       "local",
	"NAME":              "Device Simulator",
	"DESCRIPTION":       "Device Simulator",
	"HOST":              "localhost",
	"PORT":              8001,
	"ENDPOINT":          "/api/v1/device-simulator",
	"SERVER_HOST":       "localhost",
	"SERVER_PORT":       8000,
	"SERVER_ENDPOINT":   "/api/v1/master",
	"SERVER_SECRET_KEY": "secret-key",
	// DEVICE_CODE is the public identifier for the device (e.g., "DEVICE-001")
	// You can get this by creating a device via POST /api/v1/master/devices
	// The device_code will be returned in the response, or auto-generated if not provided
	"DEVICE_CODE": "DEVICE-001",
	"DEVICE_NAME": "Device 1",
	// SENSOR_ID is the ObjectID hex string of the sensor (optional, but recommended)
	// If provided, the system will infer the device from the sensor
	// You can get this by creating a sensor via POST /api/v1/master/devices/:deviceId/sensors
	"SENSOR_ID": "", // e.g., "507f1f77bcf86cd799439011"
}
