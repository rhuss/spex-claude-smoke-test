# Brainstorm: Weather ASCII CLI

**Date:** 2026-07-08
**Status:** active

## Problem Framing

Users want a quick, visually appealing way to check the weather forecast from the terminal without opening a browser or installing heavy applications. The tool should work out of the box with zero configuration, showing a 3-day forecast with large, colorful ASCII art scenes that make checking the weather enjoyable.

## Approaches Considered

### A: Single-binary, all-in-one (chosen)
- Pros: Dead simple distribution (single Go binary), easy to install via `go install`, minimal moving parts, no external dependencies
- Cons: ASCII art templates live in Go source code, updating art means recompiling

### B: Template-driven art
- Pros: ASCII art in embedded `.txt` files via Go's `embed` package, easier to edit/contribute art without touching Go code
- Cons: Slightly more complex build structure, over-engineering for a small tool

### C: Library-first with pluggable backends
- Pros: Extensible, testable via interfaces, supports multiple weather providers
- Cons: Over-engineered for a 3-day forecast CLI tool, YAGNI

## Decision

**Approach A: Single-binary, all-in-one.** Simplicity wins for a focused CLI tool. A single Go binary with embedded ASCII art and direct API calls keeps the codebase small and distribution trivial. Can evolve to template-driven art (B) later if the art library grows significantly.

## Key Requirements

- **Language:** Go, single binary with no external runtime dependencies
- **Location detection:** Auto-detect via IP geolocation API (e.g., ip-api.com), no user input needed
- **Weather data:** Open-Meteo API (free, no API key required)
- **Forecast range:** 3 days (today + 2 days)
- **ASCII art:** Large scenic art per weather condition (sun with rays, rain falling from clouds, snowflakes, etc.)
- **Color:** ANSI terminal colors (yellow sun, blue rain, white snow, gray clouds), auto-disabled when piped
- **Distribution:** Installable via `go install`

## Open Questions

- Exact set of weather conditions to illustrate (sunny, partly cloudy, cloudy, rain, heavy rain, snow, thunderstorm, fog, etc.)
- Layout preference: side-by-side days or stacked vertically?
- Which data points to show per day (min/max temp, wind speed, humidity, precipitation probability)?
- Temperature unit: Celsius default with flag for Fahrenheit, or auto-detect from locale?
- Tool name (e.g., `skycli`, `weath`, `forecast`, `skycast`)
