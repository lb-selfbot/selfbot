package conversion

var speedConverters = map[string]map[string]func(float64) float64{
	"m/s": {
		"m/s":  func(val float64) float64 { return val },
		"km/h": func(val float64) float64 { return val * 3.6 },
		"mph":  func(val float64) float64 { return val * 2.23694 },
	},
	"km/h": {
		"m/s":  func(val float64) float64 { return val / 3.6 },
		"km/h": func(val float64) float64 { return val },
		"mph":  func(val float64) float64 { return val * 0.621371 },
	},
	"mph": {
		"m/s":  func(val float64) float64 { return val / 2.23694 },
		"km/h": func(val float64) float64 { return val / 0.621371 },
		"mph":  func(val float64) float64 { return val },
	},
}

func ConvertSpeed(from string, to string, val float64) (float64, error) {
	if _, ok := speedConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := speedConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return speedConverters[from][to](val), nil
}
