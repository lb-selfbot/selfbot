package conversion

var lengthConverters = map[string]map[string]func(float64) float64{
	"m": {
		"ft": func(val float64) float64 { return val * 3.2808 },
		"in": func(val float64) float64 { return val * 39.3701 },
		"m":  func(val float64) float64 { return val },
		"km": func(val float64) float64 { return val / 1000 },
		"mi": func(val float64) float64 { return val * 0.000621371 },
	},
	"km": {
		"m":  func(val float64) float64 { return val * 1000 },
		"ft": func(val float64) float64 { return val * 3280.8 },
		"in": func(val float64) float64 { return val * 39370.1 },
		"km": func(val float64) float64 { return val },
		"mi": func(val float64) float64 { return val / 1.60934 },
	},
	"ft": {
		"m":  func(val float64) float64 { return val / 3.2808 },
		"in": func(val float64) float64 { return val * 12 },
		"ft": func(val float64) float64 { return val },
		"km": func(val float64) float64 { return val / 3280.8 },
		"mi": func(val float64) float64 { return val / 5280 },
	},
	"in": {
		"m":  func(val float64) float64 { return val / 39.3701 },
		"ft": func(val float64) float64 { return val / 12 },
		"in": func(val float64) float64 { return val },
		"km": func(val float64) float64 { return val / 39370.1 },
		"mi": func(val float64) float64 { return val / 63360 },
	},
	"mi": {
		"m":  func(val float64) float64 { return val * 1609.34 },
		"ft": func(val float64) float64 { return val * 5280 },
		"in": func(val float64) float64 { return val * 63360 },
		"km": func(val float64) float64 { return val * 1.60934 },
		"mi": func(val float64) float64 { return val },
	},
}

func ConvertLength(from string, to string, val float64) (float64, error) {
	if _, ok := lengthConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := lengthConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return lengthConverters[from][to](val), nil
}
