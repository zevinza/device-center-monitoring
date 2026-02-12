package config

var Environment = map[string]any{
	"environment":       "local",
	"name":              "Device Simulator",
	"description":       "Device Simulator",
	"host":              "localhost",
	"port":              8001,
	"endpoint":          "/api/v1/device-simulator",
	"server_host":       "localhost",
	"server_port":       8000,
	"server_endpoint":   "/api/v1/master",
	"server_secret_key": "Fr46VTqmt3j7AjT0hDa",
	"device_code":       "DEVICE-001",
	"device_name":       "Device 1",
	"sensor_id":         "698daf2a2788888746ce1b4d",
}
