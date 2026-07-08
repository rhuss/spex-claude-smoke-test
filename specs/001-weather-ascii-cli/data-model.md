# Data Model: Weather ASCII CLI

## Entities

### Location

Represents the user's detected geographic position.

| Field | Type | Description |
|-------|------|-------------|
| City | string | Human-readable city name for display |
| Lat | float64 | Latitude coordinate |
| Lon | float64 | Longitude coordinate |
| Timezone | string | IANA timezone identifier (e.g., "Europe/Berlin") |

**Source**: IP geolocation API response
**Lifecycle**: Created once per invocation, read-only after creation

### DayForecast

A single day's weather data. Three instances per invocation.

| Field | Type | Description |
|-------|------|-------------|
| Date | string | ISO date (YYYY-MM-DD) |
| Condition | WeatherCondition | Classified weather state |
| TempMax | float64 | Maximum temperature in Celsius |
| TempMin | float64 | Minimum temperature in Celsius |
| WindSpeed | float64 | Maximum wind speed in km/h |
| PrecipProb | int | Precipitation probability as percentage (0-100) |

**Source**: Open-Meteo API daily forecast response
**Lifecycle**: Created during API response parsing, read-only after creation

### WeatherCondition

Enumerated weather state that maps to ASCII art and color scheme.

| Value | WMO Codes | Description |
|-------|-----------|-------------|
| Sunny | 0, 1 | Clear sky, mainly clear |
| PartlyCloudy | 2 | Partly cloudy |
| Cloudy | 3 | Overcast |
| Fog | 45, 48 | Fog, depositing rime fog |
| Rain | 51, 53, 55, 56, 57, 61, 63, 80, 81 | Drizzle, light/moderate rain, showers |
| HeavyRain | 65, 66, 67, 82 | Heavy rain, freezing rain, violent showers |
| Snow | 71, 73, 75, 77, 85, 86 | All snow variants |
| Thunderstorm | 95, 96, 99 | Thunderstorm with/without hail |

**Source**: Derived from WMO weather code in Open-Meteo response
**Lifecycle**: Determined during forecast parsing, immutable

## Relationships

```
Location 1──fetches──▶ 3 DayForecast
DayForecast 1──has──▶ 1 WeatherCondition
WeatherCondition 1──maps-to──▶ 1 ASCII Art Scene
WeatherCondition 1──maps-to──▶ 1 Color Scheme
```

## Data Flow

1. User runs `skyweather`
2. HTTP GET to ip-api.com → Location
3. HTTP GET to Open-Meteo with Location.Lat, Location.Lon → 3x DayForecast
4. Each DayForecast.Condition → ASCII art + color selection
5. Render to stdout
