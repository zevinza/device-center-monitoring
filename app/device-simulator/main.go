package main

import (
	"api/app/device-simulator/config"
	"api/app/device-simulator/model"
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

func createRequestBody() model.SensorIngestRequest {
	now := time.Now()
	req := model.SensorIngestRequest{
		SensorID:  viper.GetString("SENSOR_ID"),
		Value:     fmt.Sprintf("%f", rand.Float64()*100),
		Timestamp: &now,
	}

	return req
}

func sendRequest(baseURL, endpoint string, body model.SensorIngestRequest) error {
	url := baseURL + endpoint

	// Marshal request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		return err
	}

	log.Println(string(jsonBody))

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
