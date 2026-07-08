// Package art provides ASCII art scenes for each weather condition.
package art

import (
	"github.com/rhuss/skyweather/internal/condition"
)

// Scene returns the ASCII art lines for a given weather condition.
// Each scene is at least 5 lines tall as required by FR-003.
func Scene(c condition.WeatherCondition) []string {
	switch c {
	case condition.Sunny:
		return sunny()
	case condition.PartlyCloudy:
		return partlyCloudy()
	case condition.Cloudy:
		return cloudy()
	case condition.Fog:
		return fog()
	case condition.Rain:
		return rain()
	case condition.HeavyRain:
		return heavyRain()
	case condition.Snow:
		return snow()
	case condition.Thunderstorm:
		return thunderstorm()
	default:
		return cloudy()
	}
}

func sunny() []string {
	return []string{
		`    \   /    `,
		`     .-.     `,
		`  - (   ) -  `,
		`     '-'     `,
		`    /   \    `,
	}
}

func partlyCloudy() []string {
	return []string{
		`   \  /      `,
		` _ /''.-.    `,
		`   \_(   ).  `,
		`   /(___(__) `,
		`             `,
	}
}

func cloudy() []string {
	return []string{
		`             `,
		`     .--.    `,
		`  .-(    ).  `,
		` (___.__)__) `,
		`             `,
	}
}

func fog() []string {
	return []string{
		`             `,
		` _ - _ - _ - `,
		`  _ - _ - _  `,
		` _ - _ - _ - `,
		`  _ - _ - _  `,
	}
}

func rain() []string {
	return []string{
		`     .-.     `,
		`    (   ).   `,
		`   (___(__)  `,
		`   ' ' ' '   `,
		`  ' ' ' '    `,
	}
}

func heavyRain() []string {
	return []string{
		`     .-.     `,
		`    (   ).   `,
		`   (___(__)  `,
		`  ,' ,' ,'   `,
		` ,' ,' ,'    `,
	}
}

func snow() []string {
	return []string{
		`     .-.     `,
		`    (   ).   `,
		`   (___(__)  `,
		`   * * * *   `,
		`  * * * *    `,
	}
}

func thunderstorm() []string {
	return []string{
		`     .-.     `,
		`    (   ).   `,
		`   (___(__)  `,
		`  ,'/_,'/_   `,
		`   /   /     `,
	}
}
