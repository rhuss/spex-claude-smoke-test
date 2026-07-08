// Package condition defines weather condition types and WMO code mapping.
package condition

// WeatherCondition represents a categorized weather state.
type WeatherCondition int

const (
	Sunny        WeatherCondition = iota // Clear sky, mainly clear
	PartlyCloudy                         // Partly cloudy
	Cloudy                               // Overcast
	Fog                                  // Fog, depositing rime fog
	Rain                                 // Drizzle, light/moderate rain, showers
	HeavyRain                            // Heavy rain, freezing rain, violent showers
	Snow                                 // All snow variants
	Thunderstorm                         // Thunderstorm with/without hail
)

// FromWMOCode converts a WMO weather interpretation code to a WeatherCondition.
// Open-Meteo uses WMO codes in its daily weather_code field.
func FromWMOCode(code int) WeatherCondition {
	switch {
	case code == 0 || code == 1:
		return Sunny
	case code == 2:
		return PartlyCloudy
	case code == 3:
		return Cloudy
	case code == 45 || code == 48:
		return Fog
	case code >= 51 && code <= 57:
		// Drizzle (51, 53, 55) and freezing drizzle (56, 57)
		return Rain
	case code == 61 || code == 63:
		// Rain slight (61) and moderate (63)
		return Rain
	case code == 65:
		// Rain heavy
		return HeavyRain
	case code == 66 || code == 67:
		// Freezing rain
		return HeavyRain
	case code >= 71 && code <= 77:
		// Snow fall (71, 73, 75) and snow grains (77)
		return Snow
	case code == 80 || code == 81:
		// Rain showers slight (80) and moderate (81)
		return Rain
	case code == 82:
		// Rain showers violent
		return HeavyRain
	case code == 85 || code == 86:
		// Snow showers
		return Snow
	case code == 95:
		// Thunderstorm
		return Thunderstorm
	case code == 96 || code == 99:
		// Thunderstorm with hail
		return Thunderstorm
	default:
		return Cloudy // Safe fallback for unknown codes
	}
}

// String returns the short name of the weather condition.
func (c WeatherCondition) String() string {
	switch c {
	case Sunny:
		return "Sunny"
	case PartlyCloudy:
		return "Partly Cloudy"
	case Cloudy:
		return "Cloudy"
	case Fog:
		return "Fog"
	case Rain:
		return "Rain"
	case HeavyRain:
		return "Heavy Rain"
	case Snow:
		return "Snow"
	case Thunderstorm:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
