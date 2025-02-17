package conversion

var areaConverters = map[string]map[string]func(float64) float64{
	"square-meter": {
		"square-meter":     func(val float64) float64 { return val },
		"square-kilometer": func(val float64) float64 { return val / 1e6 },
		"square-mile":      func(val float64) float64 { return val * 0.00000038610216 },
		"square-inch":      func(val float64) float64 { return val * 1550.0031 },
		"square-foot":      func(val float64) float64 { return val * 10.763910417 },
		"square-yard":      func(val float64) float64 { return val * 1.19599005 },
		"acre":             func(val float64) float64 { return val * 0.00024710538146717 },
	},
	"square-kilometer": {
		"square-meter":     func(val float64) float64 { return val * 1e6 },
		"square-kilometer": func(val float64) float64 { return val },
		"square-mile":      func(val float64) float64 { return val * 0.3861021585424458 },
		"square-inch":      func(val float64) float64 { return val * 15500031000 },
		"square-foot":      func(val float64) float64 { return val * 10763910.417 },
		"square-yard":      func(val float64) float64 { return val * 1195990.05 },
		"acre":             func(val float64) float64 { return val * 247.10538146717 },
	},
	"square-mile": {
		"square-meter":     func(val float64) float64 { return val * 2589988.110336 },
		"square-kilometer": func(val float64) float64 { return val * 2.589988110336 },
		"square-mile":      func(val float64) float64 { return val },
		"square-inch":      func(val float64) float64 { return val * 4014489600 },
		"square-foot":      func(val float64) float64 { return val * 27878400 },
		"square-yard":      func(val float64) float64 { return val * 3097600 },
		"acre":             func(val float64) float64 { return val * 640 },
	},
	"square-inch": {
		"square-meter":     func(val float64) float64 { return val * 0.00064516 },
		"square-kilometer": func(val float64) float64 { return val * 6.4516 * 1e-10 },
		"square-mile":      func(val float64) float64 { return val * 2.49097668 * 1e-10 },
		"square-inch":      func(val float64) float64 { return val },
		"square-foot":      func(val float64) float64 { return val * 0.006944444444444444 },
		"square-yard":      func(val float64) float64 { return val * 0.0007716049382716049 },
		"acre":             func(val float64) float64 { return val * 1.5942250790722 * 1e-10 },
	},
	"square-foot": {
		"square-meter":     func(val float64) float64 { return val * 0.09290304 },
		"square-kilometer": func(val float64) float64 { return val * 9.290304 * 1e-10 },
		"square-mile":      func(val float64) float64 { return val * 3.587006427915531 * 1e-10 },
		"square-inch":      func(val float64) float64 { return val * 144 },
		"square-foot":      func(val float64) float64 { return val },
		"square-yard":      func(val float64) float64 { return val * 0.1111111111111111 },
		"acre":             func(val float64) float64 { return val * 2.29568e-5 },
	},
	"square-yard": {
		"square-meter":     func(val float64) float64 { return val * 0.83612736 },
		"square-kilometer": func(val float64) float64 { return val * 8.3612736 * 1e-10 },
		"square-mile":      func(val float64) float64 { return val * 3.2283057851239645 * 1e-10 },
		"square-inch":      func(val float64) float64 { return val * 1296 },
		"square-foot":      func(val float64) float64 { return val * 9 },
		"square-yard":      func(val float64) float64 { return val },
		"acre":             func(val float64) float64 { return val * 0.00020661157 },
	},
	"acre": {
		"square-meter":     func(val float64) float64 { return val * 4046.8564224 },
		"square-kilometer": func(val float64) float64 { return val * 0.0040468564224 },
		"square-mile":      func(val float64) float64 { return val * 0.0015625 },
		"square-inch":      func(val float64) float64 { return val * 6272640 },
		"square-foot":      func(val float64) float64 { return val * 43560 },
		"square-yard":      func(val float64) float64 { return val * 4840 },
		"acre":             func(val float64) float64 { return val },
	},
}

func ConvertArea(from string, to string, val float64) (float64, error) {
	if _, ok := areaConverters[from]; !ok {
		return 0, ErrUnknownUnit
	}
	if _, ok := areaConverters[from][to]; !ok {
		return 0, ErrUnknownUnit
	}
	return areaConverters[from][to](val), nil
}
