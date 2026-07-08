# Implementation Plan: Weather ASCII CLI

**Branch**: `001-weather-ascii-cli` | **Date**: 2026-07-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/001-weather-ascii-cli/spec.md`

## Summary

Build a Go CLI tool that auto-detects the user's location via IP geolocation (ip-api.com), fetches a 3-day weather forecast from Open-Meteo, and renders each day as a large, colorful ASCII art scene in the terminal. The tool is distributed as a single binary with zero configuration.

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: Standard library plus `golang.org/x/term` for TTY and terminal size detection. No other third-party dependencies.
**Storage**: N/A (no persistent storage, all data fetched at runtime)
**Testing**: `go test ./...`
**Target Platform**: Cross-platform (Linux, macOS, Windows terminals with ANSI support)
**Project Type**: CLI tool
**Performance Goals**: Total execution time < 3 seconds including network requests
**Constraints**: 5-second timeout per API call, single binary, no API keys required
**Scale/Scope**: Single-user CLI tool, 2 HTTP API calls per invocation

## Global Constraints

These constraints apply to every task and are inherited implicitly:

- **Go version**: 1.22+ required
- **Dependencies**: Standard library plus `golang.org/x/term` only. No other third-party packages.
- **Network timeouts**: 5 seconds per HTTP API call. No automatic retries.
- **Distribution**: Single binary via `go install`, no external runtime dependencies (FR-010)
- **Temperature unit**: Celsius only (FR-008)
- **API keys**: None required. All APIs must be free and keyless (FR-002)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is template-only (no project-specific principles defined). No gates to enforce.

## Project Structure

### Documentation (this feature)

```text
specs/001-weather-ascii-cli/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI interface contract)
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

### Source Code (repository root)

```text
cmd/
└── skyweather/
    └── main.go          # Entry point, argument parsing, orchestration

internal/
├── geo/
│   └── geo.go           # IP geolocation client (ip-api.com)
├── weather/
│   └── weather.go       # Open-Meteo API client, response parsing
├── art/
│   └── art.go           # ASCII art scenes per weather condition
├── render/
│   └── render.go        # Layout engine (side-by-side / stacked), color output
└── condition/
    └── condition.go     # Weather condition classification (WMO code mapping)

go.mod
go.sum
```

**Structure Decision**: Standard Go project layout with `cmd/` for the binary entry point and `internal/` for private packages. Each concern (geolocation, weather API, ASCII art, rendering) is a separate package for testability. No `pkg/` directory since this is not a library.

## Complexity Tracking

No constitution violations. No complexity justifications needed.
