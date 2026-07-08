// Package geo provides IP-based geolocation using ip-api.com.
package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Location represents the user's detected geographic position.
type Location struct {
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
}

// apiResponse is the raw JSON response from ip-api.com.
type apiResponse struct {
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
}

// FetchLocation detects the user's location via IP geolocation.
// It uses ip-api.com's free tier (HTTP, no API key required).
// Network requests time out after 5 seconds.
func FetchLocation() (Location, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://ip-api.com/json/")
	if err != nil {
		return Location{}, fmt.Errorf("could not determine location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("geolocation service returned status %d", resp.StatusCode)
	}

	var result apiResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return Location{}, fmt.Errorf("could not parse geolocation response: %w", err)
	}

	if result.Status != "success" {
		return Location{}, fmt.Errorf("geolocation failed: %s", result.Message)
	}

	return Location{
		City:     result.City,
		Country:  result.Country,
		Lat:      result.Lat,
		Lon:      result.Lon,
		Timezone: result.Timezone,
	}, nil
}
