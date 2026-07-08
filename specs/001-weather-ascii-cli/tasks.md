# Tasks: Weather ASCII CLI (skyweather)

**Input**: Design documents from `specs/001-weather-ascii-cli/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: Project initialization and Go module structure

- [x] T001 Initialize Go module, create project directory structure per plan (`go.mod`, `cmd/skyweather/`, `internal/geo/`, `internal/weather/`, `internal/art/`, `internal/render/`, `internal/condition/`), and add `golang.org/x/term` dependency for TTY and terminal size detection

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and utilities that all user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 [P] Define Location struct and geolocation client in `internal/geo/geo.go` (HTTP GET to ip-api.com, JSON parsing, 5-second timeout, error handling)
  **Interfaces**: `type Location struct { City string; Lat, Lon float64; Timezone string }` | `func FetchLocation() (Location, error)`
- [x] T004 [P] Define WeatherCondition type and WMO code mapping in `internal/condition/condition.go` (enum-like constants for 8 conditions, `FromWMOCode()` function, `String()` and `Description()` methods)
  **Interfaces**: `type WeatherCondition int` with constants `Sunny, PartlyCloudy, Cloudy, Fog, Rain, HeavyRain, Snow, Thunderstorm` | `func FromWMOCode(code int) WeatherCondition` | `func (c WeatherCondition) String() string` | `func (c WeatherCondition) Description() string`
- [x] T005 [P] Define DayForecast struct and Open-Meteo client in `internal/weather/weather.go` (HTTP GET with lat/lon params, JSON parsing, 3-day daily forecast, 5-second timeout, error handling)
  **Interfaces**: `type DayForecast struct { Date string; Condition WeatherCondition; TempMax, TempMin, WindSpeed float64; PrecipProb int }` | `func FetchForecast(lat, lon float64) ([]DayForecast, error)`

**Checkpoint**: Foundation ready - all API clients and data types available

---

## Phase 3: User Story 1 - Quick Weather Check (Priority: P1) MVP

**Goal**: Run `skyweather` with no arguments, see a 3-day forecast with ASCII art for the current location

**Independent Test**: Run `go run ./cmd/skyweather/` and verify it prints location header, 3 days of weather data with ASCII art, and exits 0

### Implementation for User Story 1

- [x] T006 [P] [US1] Create ASCII art scenes for all 8 weather conditions in `internal/art/art.go` (each scene at least 5 lines tall: sunny, partly_cloudy, cloudy, fog, rain, heavy_rain, snow, thunderstorm)
- [x] T007 [P] [US1] Implement color scheme mapping in `internal/render/color.go` (ANSI color codes per condition, color-wrapping helper function, Reset code)
- [x] T008 [US1] Implement side-by-side layout renderer in `internal/render/render.go` (takes 3 DayForecast + Location, renders header with city name, 3 columns with ASCII art + data, uses terminal width from `golang.org/x/term`)
- [x] T009 [US1] Implement stacked layout fallback in `internal/render/render.go` (vertical layout for terminals < 80 columns wide, same data per day but stacked vertically)
- [x] T010 [US1] Implement main entry point in `cmd/skyweather/main.go` (orchestrate: geo lookup, weather fetch, condition mapping, render to stdout, error handling with stderr messages and exit code 1)

**Checkpoint**: User Story 1 complete - tool runs end-to-end, shows colorful 3-day forecast

---

## Phase 4: User Story 2 - Colorful Terminal Experience (Priority: P1)

**Goal**: ANSI colors applied correctly per condition, automatically disabled when piped

**Independent Test**: Run tool in terminal (colors visible), then pipe to `cat` (no ANSI codes in output)

### Implementation for User Story 2

- [x] T011 [US2] Implement TTY detection and color toggle in `internal/render/color.go` (use `golang.org/x/term.IsTerminal()` on stdout fd, expose `ColorEnabled` flag, conditionally wrap text in ANSI codes)
- [x] T012 [US2] Update `internal/render/render.go` to pass color-enabled flag through to all art and data rendering, ensuring piped output is plain text

**Checkpoint**: Colors work in terminal, clean text when piped

---

## Phase 5: User Story 3 - Graceful Error Handling (Priority: P2)

**Goal**: Clear, human-readable error messages for all failure modes

**Independent Test**: Disconnect internet and run tool, verify error message on stderr and exit code 1

### Implementation for User Story 3

- [x] T013 [US3] Add specific error messages in `cmd/skyweather/main.go` for each failure mode: geolocation unreachable ("Could not determine location"), weather API unreachable ("Could not fetch weather data"), no data for location ("No weather data available for detected location"), timeout ("Request timed out after 5 seconds")
- [x] T014 [US3] Handle partial API responses in `internal/weather/weather.go` (fewer than 3 days returned) and display available days with "Forecast data limited to N day(s)" note in `internal/render/render.go`

**Checkpoint**: All error scenarios produce clear messages, partial data handled gracefully

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final quality improvements

- [x] T015 Add `go vet` and `gofmt` checks, ensure all files pass. Verify `go install ./cmd/skyweather/` succeeds (SC-004). Time a full run to confirm < 3 seconds on a typical connection (SC-001).
- [x] T016 Write basic unit tests for WMO code mapping in `internal/condition/condition_test.go`
- [x] T017 [P] Write unit test for TTY detection logic in `internal/render/color_test.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (T003, T004, T005)
- **User Story 2 (Phase 4)**: Depends on User Story 1 (modifies render pipeline)
- **User Story 3 (Phase 5)**: Depends on User Story 1 (modifies main and render)
- **Polish (Phase 6)**: Depends on all user stories complete

### Within Each User Story

- T006 and T007 can run in parallel (different files)
- T008 depends on T006 and T007 (uses art and color)
- T009 depends on T008 (extends the renderer)
- T010 depends on all prior US1 tasks (orchestrator)

### Parallel Opportunities

- T003, T004, T005 can all run in parallel (different packages)
- T006, T007 can run in parallel (different files)
- T016, T017 can run in parallel (different test files)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T003-T005)
3. Complete Phase 3: User Story 1 (T006-T010)
4. **STOP and VALIDATE**: Run `go run ./cmd/skyweather/` - should show colored forecast
5. Working MVP with core functionality

### Incremental Delivery

1. Setup + Foundational -> API clients ready
2. User Story 1 -> Colored 3-day forecast works (MVP)
3. User Story 2 -> Pipe-safe output
4. User Story 3 -> Robust error handling
5. Polish -> Tests, formatting, code quality

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently testable after completion
- Commit after each task or logical group
- Total tasks: 16
- US1: 5 tasks, US2: 2 tasks, US3: 2 tasks, Setup: 1 task, Foundation: 3 tasks, Polish: 3 tasks
