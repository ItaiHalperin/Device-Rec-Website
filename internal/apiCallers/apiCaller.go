package apiCallers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ItaiHalperin/Device-Rec-API/dataTypes"
	"io"
	"log"
	"net/http"
)

const (
	apiHost        = "localhost:8080"
	topDevicesPath = "/api/v1/top-devices"
)

type APIResponse struct {
	Devices []dataTypes.Device `json:"devices"`
}

func GetDevices(filters *dataTypes.Filters) ([]dataTypes.Device, error) {
	// Convert the Filters object to JSON
	jsonData, err := json.Marshal(filters)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return nil, err
	}

	// Upload the JSON to the API endpoint
	url := fmt.Sprintf("http://%s%s", apiHost, topDevicesPath)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("WARNING: Failed to close HTML reader: %v", err)
		}
	}(resp.Body)

	var responseData APIResponse
	if err = json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Printf("failed to decode json from api response: %v", err)
		return nil, err
	}

	return responseData.Devices, nil
}
