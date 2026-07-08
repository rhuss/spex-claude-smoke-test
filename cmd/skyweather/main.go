// skyweather displays a colorful 3-day weather forecast in the terminal.
// It auto-detects user location via IP geolocation and fetches data from Open-Meteo.
package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/rhuss/skyweather/internal/geo"
	"github.com/rhuss/skyweather/internal/render"
	"github.com/rhuss/skyweather/internal/weather"
)

func main() {
	render.DetectTTY()

	loc, err := geo.FetchLocation()
	if err != nil {
		printError(err, "Could not determine location.")
		os.Exit(1)
	}

	forecasts, err := weather.FetchForecast(loc.Lat, loc.Lon)
	if err != nil {
		if errors.Is(err, weather.ErrNoData) {
			fmt.Fprintln(os.Stderr, "Error: No weather data available for detected location.")
		} else {
			printError(err, "Could not fetch weather data.")
		}
		os.Exit(1)
	}

	render.Forecast(os.Stdout, loc, forecasts)
}

// printError checks for timeout errors and prints the appropriate message.
func printError(err error, fallback string) {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Timeout() {
		fmt.Fprintln(os.Stderr, "Error: Request timed out after 5 seconds.")
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", fallback)
}
