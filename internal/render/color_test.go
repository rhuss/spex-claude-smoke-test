package render

import (
	"testing"

	"github.com/rhuss/skyweather/internal/condition"
)

func TestColorizeWithColorEnabled(t *testing.T) {
	ColorEnabled = true
	defer func() { ColorEnabled = true }()

	result := Colorize("hello", Yellow)
	expected := Yellow + "hello" + Reset
	if result != expected {
		t.Errorf("Colorize with color enabled = %q, want %q", result, expected)
	}
}

func TestColorizeWithColorDisabled(t *testing.T) {
	ColorEnabled = false
	defer func() { ColorEnabled = true }()

	result := Colorize("hello", Yellow)
	if result != "hello" {
		t.Errorf("Colorize with color disabled = %q, want %q", result, "hello")
	}
}

func TestColorizeWithEmptyColor(t *testing.T) {
	ColorEnabled = true
	defer func() { ColorEnabled = true }()

	result := Colorize("hello", "")
	if result != "hello" {
		t.Errorf("Colorize with empty color = %q, want %q", result, "hello")
	}
}

func TestColorForConditionEnabled(t *testing.T) {
	ColorEnabled = true
	defer func() { ColorEnabled = true }()

	tests := []struct {
		condition condition.WeatherCondition
		expected  string
	}{
		{condition.Sunny, Yellow},
		{condition.PartlyCloudy, Yellow},
		{condition.Cloudy, Gray},
		{condition.Fog, Gray},
		{condition.Rain, Blue},
		{condition.HeavyRain, BrightBlue},
		{condition.Snow, BrightWhite},
		{condition.Thunderstorm, BrightYellow},
	}

	for _, tt := range tests {
		t.Run(tt.condition.String(), func(t *testing.T) {
			got := ColorForCondition(tt.condition)
			if got != tt.expected {
				t.Errorf("ColorForCondition(%v) = %q, want %q", tt.condition, got, tt.expected)
			}
		})
	}
}

func TestColorForConditionDisabled(t *testing.T) {
	ColorEnabled = false
	defer func() { ColorEnabled = true }()

	conditions := []condition.WeatherCondition{
		condition.Sunny,
		condition.PartlyCloudy,
		condition.Cloudy,
		condition.Fog,
		condition.Rain,
		condition.HeavyRain,
		condition.Snow,
		condition.Thunderstorm,
	}

	for _, c := range conditions {
		t.Run(c.String(), func(t *testing.T) {
			got := ColorForCondition(c)
			if got != "" {
				t.Errorf("ColorForCondition(%v) with colors disabled = %q, want empty string", c, got)
			}
		})
	}
}
