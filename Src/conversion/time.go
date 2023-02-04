package conversion

var timeConverters = map[string]map[string]func(float64) float64{
	"s": {
		"s":  func(value float64) float64 { return value },
		"m":  func(value float64) float64 { return value / 60 },
		"h":  func(value float64) float64 { return value / 3600 },
		"d":  func(value float64) float64 { return value / 86400 },
		"w":  func(value float64) float64 { return value / 604800 },
		"mo": func(value float64) float64 { return value / 2628000 },
		"y":  func(value float64) float64 { return value / 31536000 },
	},
	"m": {
		"s":  func(value float64) float64 { return value * 60 },
		"m":  func(value float64) float64 { return value },
		"h":  func(value float64) float64 { return value / 60 },
		"d":  func(value float64) float64 { return value / 1440 },
		"w":  func(value float64) float64 { return value / 10080 },
		"mo": func(value float64) float64 { return value / 43800 },
		"y":  func(value float64) float64 { return value / 525600 },
	},
	"h": {
		"s":  func(value float64) float64 { return value * 3600 },
		"m":  func(value float64) float64 { return value * 60 },
		"h":  func(value float64) float64 { return value },
		"d":  func(value float64) float64 { return value / 24 },
		"w":  func(value float64) float64 { return value / 168 },
		"mo": func(value float64) float64 { return value / 730.001 },
		"y":  func(value float64) float64 { return value / 8760 },
	},
	"d": {
		"s":  func(value float64) float64 { return value * 86400 },
		"m":  func(value float64) float64 { return value * 1440 },
		"h":  func(value float64) float64 { return value * 24 },
		"d":  func(value float64) float64 { return value },
		"w":  func(value float64) float64 { return value / 7 },
		"mo": func(value float64) float64 { return value / 30.4167 },
		"y":  func(value float64) float64 { return value / 365.25 },
	},
	"w": {
		"s":  func(value float64) float64 { return value * 604800 },
		"m":  func(value float64) float64 { return value * 10080 },
		"h":  func(value float64) float64 { return value * 168 },
		"d":  func(value float64) float64 { return value * 7 },
		"w":  func(value float64) float64 { return value },
		"mo": func(value float64) float64 { return value / 4.34524 },
		"y":  func(value float64) float64 { return value / 52.1775 },
	},
	"mo": {
		"s":  func(value float64) float64 { return value * 2628000 },
		"m":  func(value float64) float64 { return value * 43800 },
		"h":  func(value float64) float64 { return value * 730.001 },
		"d":  func(value float64) float64 { return value * 30.4167 },
		"w":  func(value float64) float64 { return value * 4.34524 },
		"mo": func(value float64) float64 { return value },
		"y":  func(value float64) float64 { return value / 12 },
	},
	"y": {
		"s":  func(value float64) float64 { return value * 31536000 },
		"m":  func(value float64) float64 { return value * 525600 },
		"h":  func(value float64) float64 { return value * 8760 },
		"d":  func(value float64) float64 { return value * 365.25 },
		"w":  func(value float64) float64 { return value * 52.1775 },
		"mo": func(value float64) float64 { return value * 12 },
		"y":  func(value float64) float64 { return value },
	},
}

func ConvertTime(from string, to string, val float64) (float64, error) {
	if _, ok := timeConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := timeConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return timeConverters[from][to](val), nil
}
