package condition

import "testing"

func TestFromWMOCode(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected WeatherCondition
	}{
		// Sunny (clear sky, mainly clear)
		{"clear sky", 0, Sunny},
		{"mainly clear", 1, Sunny},

		// Partly cloudy
		{"partly cloudy", 2, PartlyCloudy},

		// Cloudy (overcast)
		{"overcast", 3, Cloudy},

		// Fog
		{"fog", 45, Fog},
		{"rime fog", 48, Fog},

		// Rain (drizzle, light/moderate rain, showers)
		{"drizzle light", 51, Rain},
		{"drizzle moderate", 53, Rain},
		{"drizzle dense", 55, Rain},
		{"freezing drizzle light", 56, Rain},
		{"freezing drizzle dense", 57, Rain},
		{"rain slight", 61, Rain},
		{"rain moderate", 63, Rain},
		{"rain showers slight", 80, Rain},
		{"rain showers moderate", 81, Rain},

		// Heavy rain
		{"rain heavy", 65, HeavyRain},
		{"freezing rain light", 66, HeavyRain},
		{"freezing rain heavy", 67, HeavyRain},
		{"rain showers violent", 82, HeavyRain},

		// Snow
		{"snow slight", 71, Snow},
		{"snow moderate", 73, Snow},
		{"snow heavy", 75, Snow},
		{"snow grains", 77, Snow},
		{"snow showers slight", 85, Snow},
		{"snow showers heavy", 86, Snow},

		// Thunderstorm
		{"thunderstorm", 95, Thunderstorm},
		{"thunderstorm with hail slight", 96, Thunderstorm},
		{"thunderstorm with hail heavy", 99, Thunderstorm},

		// Unknown code falls back to Cloudy
		{"unknown code", 999, Cloudy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromWMOCode(tt.code)
			if got != tt.expected {
				t.Errorf("FromWMOCode(%d) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestWeatherConditionString(t *testing.T) {
	tests := []struct {
		condition WeatherCondition
		expected  string
	}{
		{Sunny, "Sunny"},
		{PartlyCloudy, "Partly Cloudy"},
		{Cloudy, "Cloudy"},
		{Fog, "Fog"},
		{Rain, "Rain"},
		{HeavyRain, "Heavy Rain"},
		{Snow, "Snow"},
		{Thunderstorm, "Thunderstorm"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.condition.String()
			if got != tt.expected {
				t.Errorf("%v.String() = %q, want %q", tt.condition, got, tt.expected)
			}
		})
	}
}
