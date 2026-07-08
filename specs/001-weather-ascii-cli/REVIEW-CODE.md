# Code Review Report

**Feature**: 001-weather-ascii-cli
**Date**: 2026-07-08
**Reviewer**: Claude Code (automated)

## Spec Compliance Check

**Score: 14/14 (100%)**

| Requirement | Status | Evidence |
|-------------|--------|----------|
| FR-001: Auto-detect location via IP geolocation | PASS | `geo.FetchLocation()` calls `ip-api.com/json/` with no user input required |
| FR-002: 3-day forecast from free API (no key) | PASS | `weather.FetchForecast()` calls Open-Meteo with `forecast_days=3`, no API key |
| FR-003: ASCII art at least 5 lines tall per day | PASS | `art.Scene()` returns exactly 5 lines for all 8 conditions |
| FR-004: 8+ distinct weather conditions with unique art | PASS | 8 conditions: sunny, partly cloudy, cloudy, fog, rain, heavy rain, snow, thunderstorm |
| FR-005: ANSI colors matching conditions | PASS | `ColorForCondition()`: yellow/sun, blue/rain, white/snow, gray/clouds |
| FR-006: TTY detection, suppress colors when piped | PASS | `DetectTTY()` uses `term.IsTerminal()`, sets `ColorEnabled` |
| FR-007: Day name, date, condition, temps, wind, precip | PASS | Both `buildDayColumn` and `renderStackedDay` display all required fields |
| FR-008: Temperatures in Celsius | PASS | Open-Meteo defaults to Celsius; display uses `C` suffix |
| FR-009: Detected city name in header | PASS | `Forecast()` prints `"Weather for <City>, <Country>"` |
| FR-010: Single binary, no runtime dependencies | PASS | Go compiles to static binary; only stdlib + `golang.org/x/term` |
| FR-011: Meaningful error messages on network failure | PASS | `printError()` provides human-readable messages for all error paths |
| FR-012: Side-by-side layout at >= 80 columns | PASS | `Forecast()` checks `width >= 80`, calls `renderSideBySide()` |
| FR-013: Stacked layout at < 80 columns | PASS | Else branch calls `renderStacked()` |
| FR-014: 5-second timeout, error on timeout | PASS | Both HTTP clients set `Timeout: 5 * time.Second`; `printError()` detects timeouts |

### Edge Cases from Spec

| Edge Case | Status | Evidence |
|-----------|--------|----------|
| No weather coverage for location | PASS | Returns `"no weather data available for detected location"` error |
| Narrow terminal (<80 cols) fallback | PASS | Stacked layout triggered when `width < 80` |
| VPN location behavior | PASS | Uses whatever ip-api.com returns (documented as acceptable) |
| Partial data (<3 days) | PASS | Shows available days + `"Forecast data limited to N day(s)."` note |
| 5-second timeout, no retries | PASS | Timeout set on both clients; no retry logic; non-zero exit on error |

## CodeRabbit Review

CodeRabbit found 1 major finding:

- **[FIXED] `internal/weather/weather.go:71-90`**: `WeatherCode[i]` indexed without bounds check, unlike other fields. Added bounds guard consistent with `TempMax`, `TempMin`, `WindMax`, and `PrecipProb`.

## Deep Review Report

### 1. Correctness Review

**Severity: 1 Critical (fixed), 1 Minor**

- **[FIXED] Panic on partial API response** (`internal/weather/weather.go:76`): `days.WeatherCode[i]` was accessed without bounds checking, while all other daily fields (`TempMax`, `TempMin`, `WindMax`, `PrecipProb`) had proper `i < len(...)` guards. If the Open-Meteo API returned a `time` array longer than `weather_code`, the program would panic with an index-out-of-range error. Fixed by adding a bounds check consistent with the other fields.

- **[Minor] Default terminal width when piped** (`internal/render/render.go:48-55`): `terminalWidth()` defaults to 80 when `term.GetSize()` fails (e.g., piped output). This means piped output always uses the side-by-side layout, which is reasonable behavior but could produce wide output that wraps in narrow terminals when re-displayed. No fix needed; behavior is acceptable per spec.

### 2. Architecture Review

**Severity: No issues**

The codebase follows clean Go architecture:

- **Package separation**: Clear boundaries between `geo` (location), `weather` (API client), `condition` (domain types), `art` (ASCII art), and `render` (terminal output).
- **Dependency direction**: Packages depend inward on `condition` (the domain type); `render` orchestrates the others. No circular dependencies.
- **Single responsibility**: Each package owns one concern. The `condition` package cleanly separates WMO code mapping from both the API client and the renderer.
- **Appropriate use of interfaces**: `render.Forecast()` accepts `io.Writer` for testability.

### 3. Security Review

**Severity: 1 Low**

- **[Low] Unencrypted geolocation request** (`internal/geo/geo.go:39`): Uses `http://` (not `https://`) for `ip-api.com/json/`. This is because ip-api.com's free tier only supports HTTP. The response contains city-level location data derived from the user's IP address. For a local CLI weather tool, the risk is low (the data is not sensitive beyond what the IP itself reveals), but users on untrusted networks could have their approximate location observed in transit. No fix applied; this is a limitation of the free API tier.
- No user input is passed unsanitized to any API call. Latitude/longitude are `float64` values from JSON decoding, formatted with `%.4f`.
- No secrets, credentials, or API keys are stored or transmitted.

### 4. Production Readiness Review

**Severity: 1 Important, 1 Minor**

- **[Important] Compiled binary tracked in git** (`skyweather`, 8.2MB): The compiled binary is checked into the repository. This bloats the repo, makes diffs meaningless for the binary, and the binary may not work on different OS/architectures. The `.gitignore` file exists but is untracked. Recommendation: add `skyweather` to `.gitignore` and remove it from tracking.

- **[Minor] No version/help flags**: The CLI has no `--version` or `--help` flag. While not required by the spec, these are standard for production CLI tools. Not blocking for v1.

### 5. Test Coverage Review

**Severity: 1 Important**

- **[Important] Limited test coverage**: Only `internal/condition` and `internal/render` (color functions) have unit tests. The following packages have no tests:
  - `internal/art` - No tests verifying scene line counts or condition coverage
  - `internal/geo` - No tests (would require HTTP mocking)
  - `internal/weather` - No tests (would require HTTP mocking)
  - `internal/render` (rendering logic) - No tests for `Forecast()`, `renderSideBySide()`, `renderStacked()`, or layout helpers

  The existing 71 tests cover WMO code mapping (28 cases), condition string/description methods (16 cases), and color functions (27 cases). These are the pure-function utilities. The core rendering and API interaction logic is untested.

  For v1, the tested functions cover the most error-prone mapping logic. HTTP-dependent packages would benefit from interface-based mocking in a future iteration.

### Summary

| Perspective | Critical | Important | Minor | Low |
|-------------|----------|-----------|-------|-----|
| Correctness | 1 (fixed) | 0 | 1 | 0 |
| Architecture | 0 | 0 | 0 | 0 |
| Security | 0 | 0 | 0 | 1 |
| Production | 0 | 1 | 1 | 0 |
| Tests | 0 | 1 | 0 | 0 |
| **Total** | **1 (fixed)** | **2** | **2** | **1** |

### Gate Decision

**PASS** - The one critical finding (WeatherCode bounds check) has been fixed. The remaining findings are Important or lower severity and do not block shipping:
- Binary tracked in git is a repo hygiene issue, not a runtime defect
- Limited test coverage is acceptable for v1 given the tested code covers the most error-prone mapping logic

Spec compliance: **14/14 (100%)**
All tests: **71/71 passing**
go vet: **Clean**
