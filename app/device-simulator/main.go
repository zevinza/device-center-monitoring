package main

import (
	"api/app/device-simulator/config"
	"api/utils/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func init() {
	env.Load(config.Environment)
}

func main() {
	// Build master service URL
	serverEndpoint := viper.GetString("SERVER_ENDPOINT") + "/sensors"
	baseURL := fmt.Sprintf("http://%s:%d", viper.GetString("SERVER_HOST"), viper.GetInt("SERVER_PORT"))

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	log.Printf("Starting device simulator - sending requests to %s%s every 3 seconds", baseURL, serverEndpoint)

	body := createRequestBody()
	if err := sendRequest(baseURL, serverEndpoint, body); err != nil {
		log.Printf("Error sending request: %v", err)
	}

	for range ticker.C {
		body := createRequestBody()
		if err := sendRequest(baseURL, serverEndpoint, body); err != nil {
			log.Printf("Error sending request: %v", err)
		}
	}
}

type BodyRequest struct {
	DeviceCode *string `json:"device_code,omitempty"` // Optional if sensor_id is provided
	SensorID   *string `json:"sensor_id,omitempty"`   // Recommended: system infers device from sensor
	Unit       string  `json:"unit"`
	Value      float64 `json:"value"`
}

func createRequestBody() BodyRequest {
	req := BodyRequest{
		Unit:  "celsius",
		Value: rand.Float64() * 100,
	}

	// If SENSOR_ID is configured, use it (recommended approach)
	// Otherwise, fall back to DEVICE_CODE (backward compatibility)
	sensorID := viper.GetString("SENSOR_ID")
	if sensorID != "" {
		req.SensorID = &sensorID
	} else {
		deviceCode := viper.GetString("DEVICE_CODE")
		if deviceCode != "" {
			req.DeviceCode = &deviceCode
		}
	}

	return req
}

func sendRequest(baseURL, endpoint string, body BodyRequest) error {
	url := baseURL + endpoint

	// Marshal request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", viper.GetString("SERVER_SECRET_KEY"))

	// Send request
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}

	log.Printf("Response from master service [%d]: %s", resp.StatusCode, string(respBody))

	return nil
}
