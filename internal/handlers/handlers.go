package handlers

import (
	"deviceRecommendationWebsite/internal/apiCallers"
	"encoding/json"
	"github.com/ItaiHalperin/deviceRecsResources/dataTypes"
	"log"
	"net/http"
	"strings"
)

// Message struct to handle incoming JSON
type Message struct {
	Brands      []string           `json:"brands"`
	PriceRange  map[string]int     `json:"priceRange"`
	DisplaySize map[string]float64 `json:"displaySize"`
	RefreshRate map[string]int     `json:"refreshRate"`
}

// Response struct to send back to client
type Response struct {
	Results []dataTypes.Device `json:"results"`
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are accepted
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON message
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received choice
	log.Printf("Received choice: %s, %v, %v, %v\n", msg.Brands, msg.PriceRange, msg.RefreshRate, msg.DisplaySize)

	var capitalizedBrandNames []string
	for _, brand := range msg.Brands {
		capitalizedBrandNames = append(capitalizedBrandNames, strings.ToUpper(brand[:1])+brand[1:])
	}
	filters := dataTypes.Filters{
		Price:       dataTypes.MinMaxInt{Min: msg.PriceRange["min"], Max: msg.PriceRange["max"]},
		DisplaySize: dataTypes.MinMaxFloat{Min: msg.DisplaySize["min"], Max: msg.DisplaySize["max"]},
		RefreshRate: dataTypes.MinMaxInt{Min: msg.RefreshRate["min"], Max: msg.RefreshRate["max"]},
		Brands:      capitalizedBrandNames,
	}

	devices, err := apiCallers.GetDevices(&filters)
	if err != nil {
		log.Printf("Error getting devices: %v", err)
		http.Error(w, "Error getting devices", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send JSON response
	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		log.Printf("Error encoding devices: %v", err)
		http.Error(w, "Error encoding devices", http.StatusInternalServerError)
		return
	}
}
