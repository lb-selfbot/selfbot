package conversion

var temperatureConverters = map[string]map[string]func(float64) float64{
	"C": {
        "C": func(val float64) float64 { return val },
        "F": func(val float64) float64 { return val*9/5 + 32 },
        "K": func(val float64) float64 { return val + 273.15 },
    },
    "F": {
        "C": func(val float64) float64 { return (val - 32) * 5/9 },
        "F": func(val float64) float64 { return val },
        "K": func(val float64) float64 { return (val + 459.67) * 5/9 },
    },
    "K": {
        "C": func(val float64) float64 { return val - 273.15 },
        "F": func(val float64) float64 { return val*9/5 - 459.67 },
        "K": func(val float64) float64 { return val },
    },
}

func ConvertTemperature(from string, to string, val float64) (float64, error) {
	if _, ok := temperatureConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := temperatureConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return temperatureConverters[from][to](val), nil
}
