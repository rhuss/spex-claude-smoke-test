# Research: Weather ASCII CLI

## API Selection

### IP Geolocation: ip-api.com

**Decision**: Use ip-api.com free tier (HTTP, no key required)
**Rationale**: Free for non-commercial use, returns city, latitude, longitude in JSON. No registration needed. HTTP endpoint at `http://ip-api.com/json/` returns location for the caller's IP.
**Alternatives considered**:
- ipinfo.io: Requires API token for reliable access
- freegeoip.app: Deprecated
- ipapi.co: Rate-limited without token

**Response format**:
```json
{
  "status": "success",
  "city": "Berlin",
  "lat": 52.5200,
  "lon": 13.4050,
  "timezone": "Europe/Berlin"
}
```

**Rate limit**: 45 requests per minute (more than sufficient for a CLI tool).

### Weather Data: Open-Meteo

**Decision**: Use Open-Meteo forecast API (HTTPS, no key required)
**Rationale**: Free, open-source weather API. Returns daily forecasts with all needed data points. No registration or API key. Well-documented JSON responses.
**Alternatives considered**:
- OpenWeatherMap: Requires API key registration
- wttr.in: Returns pre-rendered ASCII (limits creative control)
- WeatherAPI.com: Requires API key

**Endpoint**: `https://api.open-meteo.com/v1/forecast`
**Parameters**:
- `latitude`, `longitude`: From geolocation
- `daily`: `weather_code,temperature_2m_max,temperature_2m_min,wind_speed_10m_max,precipitation_probability_max`
- `forecast_days`: 3
- `timezone`: auto

**Response format** (daily fields):
```json
{
  "daily": {
    "time": ["2026-07-08", "2026-07-09", "2026-07-10"],
    "weather_code": [1, 61, 71],
    "temperature_2m_max": [28.5, 22.1, 18.3],
    "temperature_2m_min": [18.2, 15.4, 12.1],
    "wind_speed_10m_max": [15.2, 22.5, 10.1],
    "precipitation_probability_max": [10, 75, 40]
  }
}
```

## WMO Weather Codes to Conditions

Open-Meteo uses WMO weather interpretation codes. Mapping to our 8 condition categories:

| WMO Code | Description | Our Category |
|----------|-------------|-------------|
| 0 | Clear sky | sunny |
| 1 | Mainly clear | sunny |
| 2 | Partly cloudy | partly_cloudy |
| 3 | Overcast | cloudy |
| 45, 48 | Fog, rime fog | fog |
| 51, 53, 55 | Drizzle | rain |
| 56, 57 | Freezing drizzle | rain |
| 61, 63 | Rain (slight, moderate) | rain |
| 65 | Rain (heavy) | heavy_rain |
| 66, 67 | Freezing rain | heavy_rain |
| 71, 73, 75, 77 | Snow | snow |
| 80, 81, 82 | Rain showers | rain / heavy_rain |
| 85, 86 | Snow showers | snow |
| 95 | Thunderstorm | thunderstorm |
| 96, 99 | Thunderstorm with hail | thunderstorm |

## TTY Detection

**Decision**: Use `os.Stdout.Fd()` with `golang.org/x/term.IsTerminal()` or check `os.Getenv("TERM")` and `os.IsNotExist` on stat of stdout.
**Rationale**: Go standard library does not have a built-in TTY check. However, we can use `golang.org/x/term` (a quasi-standard library maintained by the Go team) or implement a simple syscall-based check. Since the spec says "no external runtime dependencies" (referring to the binary), build-time dependencies like `golang.org/x/term` are acceptable.
**Decision**: Use `golang.org/x/term` for TTY detection. It's a Go team-maintained module and compiles into the binary (no runtime dependency).

## Terminal Width Detection

**Decision**: Use `golang.org/x/term.GetSize()` for terminal width detection.
**Rationale**: Same module used for TTY detection. Returns terminal dimensions. Fallback to 80 columns if detection fails (e.g., when piped).

## Tool Name

**Decision**: `skyweather`
**Rationale**: Descriptive, memorable, not taken on GitHub. The `sky` prefix hints at the visual/ASCII art nature. Short enough for frequent CLI use.
**Alternatives considered**:
- `weath`: Too abbreviated, unclear
- `skycli`: Generic "cli" suffix adds no meaning
- `forecast`: Too common, likely name conflicts
- `skycast`: Good but `skyweather` is more immediately clear

## Color Scheme

| Condition | Primary Color | ANSI Code |
|-----------|--------------|-----------|
| sunny | Yellow | `\033[33m` |
| partly_cloudy | Yellow + White | `\033[33m` / `\033[37m` |
| cloudy | Gray (bright black) | `\033[90m` |
| rain | Blue | `\033[34m` |
| heavy_rain | Bright blue | `\033[94m` |
| snow | Bright white | `\033[97m` |
| thunderstorm | Bright yellow (lightning) + blue | `\033[93m` / `\033[34m` |
| fog | Gray | `\033[90m` |
