// Package weather provides the Open-Meteo API client for fetching forecasts.
package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rhuss/skyweather/internal/condition"
)

// DayForecast holds a single day's weather data.
type DayForecast struct {
	Date       string                     // ISO date (YYYY-MM-DD)
	Condition  condition.WeatherCondition // Classified weather state
	TempMax    float64                    // Maximum temperature in Celsius
	TempMin    float64                    // Minimum temperature in Celsius
	WindSpeed  float64                    // Maximum wind speed in km/h
	PrecipProb int                        // Precipitation probability percentage (0-100)
}

// apiResponse is the raw JSON response from Open-Meteo.
type apiResponse struct {
	Daily dailyData `json:"daily"`
}

type dailyData struct {
	Time        []string  `json:"time"`
	WeatherCode []int     `json:"weather_code"`
	TempMax     []float64 `json:"temperature_2m_max"`
	TempMin     []float64 `json:"temperature_2m_min"`
	WindMax     []float64 `json:"wind_speed_10m_max"`
	PrecipProb  []int     `json:"precipitation_probability_max"`
}

// ErrNoData is returned when the weather API has no forecast data for the location.
var ErrNoData = errors.New("no weather data available for detected location")

// FetchForecast retrieves a 3-day weather forecast from Open-Meteo.
// It uses the free API (no key required). Network requests time out after 5 seconds.
func FetchForecast(lat, lon float64) ([]DayForecast, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&daily=weather_code,temperature_2m_max,temperature_2m_min,wind_speed_10m_max,precipitation_probability_max"+
			"&forecast_days=3&timezone=auto",
		lat, lon,
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather service returned status %d", resp.StatusCode)
	}

	var result apiResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not parse weather response: %w", err)
	}

	days := result.Daily
	if len(days.Time) == 0 {
		return nil, ErrNoData
	}

	forecasts := make([]DayForecast, len(days.Time))
	for i := range days.Time {
		f := DayForecast{
			Date:      days.Time[i],
			Condition: condition.Cloudy,
		}
		if i < len(days.WeatherCode) {
			f.Condition = condition.FromWMOCode(days.WeatherCode[i])
		}
		if i < len(days.TempMax) {
			f.TempMax = days.TempMax[i]
		}
		if i < len(days.TempMin) {
			f.TempMin = days.TempMin[i]
		}
		if i < len(days.WindMax) {
			f.WindSpeed = days.WindMax[i]
		}
		if i < len(days.PrecipProb) {
			f.PrecipProb = days.PrecipProb[i]
		}
		forecasts[i] = f
	}

	return forecasts, nil
}
