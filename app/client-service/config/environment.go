package config

var Environment = map[string]any{
	"environment": "local",
	"name":        "Client Service",
	"description": "Client Service - Receives sensor data from master service",
	"port":        8002,
	"host":        "localhost",
}
