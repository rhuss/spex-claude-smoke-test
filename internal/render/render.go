package render

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/rhuss/skyweather/internal/art"
	"github.com/rhuss/skyweather/internal/geo"
	"github.com/rhuss/skyweather/internal/weather"
)

// colWidth is the width of each day column in side-by-side layout.
const colWidth = 25

// Forecast renders the weather forecast to the given writer.
// It auto-detects terminal width and uses side-by-side layout (>= 80 cols)
// or stacked layout (< 80 cols).
func Forecast(w io.Writer, loc geo.Location, days []weather.DayForecast) {
	width := terminalWidth()

	// Header with city name and country
	header := fmt.Sprintf("Weather for %s", loc.City)
	if loc.Country != "" {
		header = fmt.Sprintf("Weather for %s, %s", loc.City, loc.Country)
	}
	fmt.Fprintln(w, Colorize(header, BrightWhite))
	fmt.Fprintln(w)

	if width >= 80 {
		renderSideBySide(w, days)
	} else {
		renderStacked(w, days)
	}

	// Handle partial data (fewer than 3 days)
	if len(days) > 0 && len(days) < 3 {
		fmt.Fprintf(w, "\nForecast data limited to %d day(s).\n", len(days))
	}
}

// terminalWidth returns the current terminal width, defaulting to 80 if
// detection fails (e.g., when output is piped).
func terminalWidth() int {
	fd := int(os.Stdout.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil || width <= 0 {
		return 80
	}
	return width
}

// renderSideBySide renders all days in a horizontal row.
func renderSideBySide(w io.Writer, days []weather.DayForecast) {
	n := len(days)
	if n == 0 {
		return
	}

	// Build column content for each day
	columns := make([][]string, n)
	for i, day := range days {
		columns[i] = buildDayColumn(day)
	}

	// Top border
	fmt.Fprint(w, "┌")
	for i := range n {
		fmt.Fprint(w, strings.Repeat("─", colWidth))
		if i < n-1 {
			fmt.Fprint(w, "┬")
		}
	}
	fmt.Fprintln(w, "┐")

	// Find max rows across columns
	maxRows := 0
	for _, col := range columns {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}

	// Render rows
	for row := 0; row < maxRows; row++ {
		fmt.Fprint(w, "│")
		for i, col := range columns {
			cell := ""
			if row < len(col) {
				cell = col[row]
			}
			fmt.Fprint(w, padToWidth(cell, colWidth))
			if i < n-1 {
				fmt.Fprint(w, "│")
			}
		}
		fmt.Fprintln(w, "│")
	}

	// Bottom border
	fmt.Fprint(w, "└")
	for i := range n {
		fmt.Fprint(w, strings.Repeat("─", colWidth))
		if i < n-1 {
			fmt.Fprint(w, "┴")
		}
	}
	fmt.Fprintln(w, "┘")
}

// renderStacked renders each day as a separate block.
func renderStacked(w io.Writer, days []weather.DayForecast) {
	for i, day := range days {
		if i > 0 {
			fmt.Fprintln(w)
		}
		renderStackedDay(w, day)
	}
}

// renderStackedDay renders a single day in stacked layout.
func renderStackedDay(w io.Writer, day weather.DayForecast) {
	color := ColorForCondition(day.Condition)
	dayName, dateStr := formatDate(day.Date)

	// Day header
	fmt.Fprintln(w, Colorize(fmt.Sprintf("%s %s", dayName, dateStr), BrightWhite))
	fmt.Fprintln(w)

	// ASCII art
	scene := art.Scene(day.Condition)
	for _, line := range scene {
		fmt.Fprintln(w, Colorize(line, color))
	}
	fmt.Fprintln(w)

	// Weather data
	fmt.Fprintln(w, Colorize(day.Condition.String(), color))
	fmt.Fprintf(w, "%d°C / %d°C\n", int(math.Round(day.TempMax)), int(math.Round(day.TempMin)))
	fmt.Fprintf(w, "Wind: %d km/h\n", int(math.Round(day.WindSpeed)))
	fmt.Fprintf(w, "Precip: %d%%\n", day.PrecipProb)
}

// buildDayColumn builds the content lines for a single day column.
func buildDayColumn(day weather.DayForecast) []string {
	color := ColorForCondition(day.Condition)
	dayName, dateStr := formatDate(day.Date)

	var lines []string

	// Day header (centered)
	header := fmt.Sprintf("%s %s", dayName, dateStr)
	lines = append(lines, centerText(Colorize(header, BrightWhite), colWidth))
	lines = append(lines, "")

	// ASCII art
	scene := art.Scene(day.Condition)
	for _, artLine := range scene {
		lines = append(lines, centerText(Colorize(artLine, color), colWidth))
	}
	lines = append(lines, "")

	// Weather info
	lines = append(lines, centerText(Colorize(day.Condition.String(), color), colWidth))
	lines = append(lines, centerText(fmt.Sprintf("%d°C / %d°C", int(math.Round(day.TempMax)), int(math.Round(day.TempMin))), colWidth))
	lines = append(lines, centerText(fmt.Sprintf("Wind: %d km/h", int(math.Round(day.WindSpeed))), colWidth))
	lines = append(lines, centerText(fmt.Sprintf("Precip: %d%%", day.PrecipProb), colWidth))

	return lines
}

// formatDate parses an ISO date string and returns the day name and short date.
func formatDate(dateStr string) (string, string) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "Unknown", dateStr
	}
	return t.Weekday().String(), t.Format("Jan 02")
}

// centerText centers text within the given width.
// It accounts for ANSI escape codes when calculating visible length.
func centerText(text string, width int) string {
	visible := visibleLen(text)
	if visible >= width {
		return text
	}
	leftPad := (width - visible) / 2
	rightPad := width - visible - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// padToWidth pads text to the given width, accounting for ANSI codes.
func padToWidth(text string, width int) string {
	visible := visibleLen(text)
	if visible >= width {
		return text
	}
	return text + strings.Repeat(" ", width-visible)
}

// visibleLen returns the length of text excluding ANSI escape sequences.
func visibleLen(s string) int {
	inEscape := false
	count := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		count++
	}
	return count
}
