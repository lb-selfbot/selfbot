package conversion

var massConverters = map[string]map[string]func(float64) float64{
	"kg": {
		"kg": func(val float64) float64 { return val },
		"g": func(val float64) float64 { return val * 1000 },
		"mg": func(val float64) float64 { return val * 1000000 },
		"lb": func(val float64) float64 { return val * 2.20462 },
		"oz": func(val float64) float64 { return val * 35.274 },
	},
	"g": {
		"kg": func(val float64) float64 { return val / 1000 },
		"g": func(val float64) float64 { return val },
		"mg": func(val float64) float64 { return val * 1000 },
		"lb": func(val float64) float64 { return val / 453.592 * 2.20462 },
		"oz": func(val float64) float64 { return val / 28.3495 },
	},
	"mg": {
		"kg": func(val float64) float64 { return val / 1000000 },
		"g": func(val float64) float64 { return val / 1000 },
		"mg": func(val float64) float64 { return val },
		"lb": func(val float64) float64 { return val / 453592 * 2.20462 },
		"oz": func(val float64) float64 { return val / 28349.5 },
	},
	"lb": {
		"kg": func(val float64) float64 { return val / 2.20462 },
		"g": func(val float64) float64 { return val * 453.592 },
		"mg": func(val float64) float64 { return val * 453592 },
		"lb": func(val float64) float64 { return val },
		"oz": func(val float64) float64 { return val * 16 },
	},
	"oz": {
		"kg": func(val float64) float64 { return val / 35.274 },
		"g": func(val float64) float64 { return val * 28.3495 },
		"mg": func(val float64) float64 { return val * 28349.5 },
		"lb": func(val float64) float64 { return val / 16 },
		"oz": func(val float64) float64 { return val },
	},
}

func ConvertMass(from string, to string, val float64) (float64, error) {
	if _, ok := massConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := massConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return massConverters[from][to](val), nil
}
