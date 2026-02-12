package config

var Environment = map[string]any{
	"ENVIRONMENT": "local",
	"NAME":        "Client Service",
	"DESCRIPTION": "Client Service - Receives sensor data from master service",
	"PORT":        8002,
	"HOST":        "localhost",
}
