// Package render provides terminal rendering with ANSI color support.
package render

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/rhuss/skyweather/internal/condition"
)

// ANSI color escape codes.
const (
	Reset        = "\033[0m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Gray         = "\033[90m"
	BrightBlue   = "\033[94m"
	BrightWhite  = "\033[97m"
	BrightYellow = "\033[93m"
)

// ColorEnabled controls whether ANSI codes are emitted.
// Set to false when output is piped (not a TTY).
var ColorEnabled = true

// DetectTTY checks whether stdout is connected to a terminal and sets
// ColorEnabled accordingly. When output is piped to another program or
// file, ANSI color codes are suppressed.
func DetectTTY() {
	ColorEnabled = term.IsTerminal(int(os.Stdout.Fd()))
}

// ColorForCondition returns the primary ANSI color code for a weather condition.
func ColorForCondition(c condition.WeatherCondition) string {
	if !ColorEnabled {
		return ""
	}
	switch c {
	case condition.Sunny:
		return Yellow
	case condition.PartlyCloudy:
		return Yellow
	case condition.Cloudy:
		return Gray
	case condition.Fog:
		return Gray
	case condition.Rain:
		return Blue
	case condition.HeavyRain:
		return BrightBlue
	case condition.Snow:
		return BrightWhite
	case condition.Thunderstorm:
		return BrightYellow
	default:
		return ""
	}
}

// Colorize wraps text in the given ANSI color code.
// Returns plain text when ColorEnabled is false.
func Colorize(text, color string) string {
	if !ColorEnabled || color == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}
