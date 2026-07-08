# CLI Interface Contract: skyweather

## Command

```
skyweather
```

No arguments, no flags, no configuration. Zero-config by design.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, forecast displayed |
| 1 | Error (network failure, API error, no data) |

## Standard Output

### Success (wide terminal, >= 80 columns)

```
Weather for Berlin, Germany

┌─────────────────┬─────────────────┬─────────────────┐
│   Tuesday 08    │  Wednesday 09   │   Thursday 10   │
│                 │                 │                 │
│   [ASCII ART]   │   [ASCII ART]   │   [ASCII ART]   │
│   (5+ lines)    │   (5+ lines)    │   (5+ lines)    │
│                 │                 │                 │
│   Sunny         │   Rain          │   Cloudy        │
│   22°C / 15°C   │   18°C / 12°C   │   20°C / 14°C   │
│   Wind: 12 km/h │   Wind: 25 km/h │   Wind: 8 km/h  │
│   Precip: 5%    │   Precip: 80%   │   Precip: 20%   │
└─────────────────┴─────────────────┴─────────────────┘
```

### Success (narrow terminal, < 80 columns)

Days stacked vertically, one block per day.

### Partial Data

Normal output for available days, followed by:
```
Forecast data limited to N day(s).
```

## Standard Error

All error messages go to stderr.

```
Error: Could not determine location. Check your internet connection.
Error: Could not fetch weather data. Check your internet connection.
Error: No weather data available for detected location.
Error: Request timed out after 5 seconds.
```

## Piped Output

When stdout is not a TTY, all ANSI escape codes are omitted. Plain ASCII text only.
