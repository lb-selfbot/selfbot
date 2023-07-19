package conversion

import "errors"

var AreaUnits = []string{
	"square-meter",
	"square-kilometer",
	"square-mile",
	"square-inch",
	"square-foot",
	"square-yard",
	"acre",
}

var LengthUnits = []string{
	"m",
	"ft",
	"in",
	"km",
	"mi",
}

var TimeUnits = []string{
	"s",
	"m",
	"h",
	"d",
	"w",
	"mo",
	"y",
}

var TemperatureUnits = []string{"C", "F", "K"}

var MassUnits = []string{
	"kg",
	"g",
	"mg",
	"lb",
	"oz",
}

var SpeedUnits = []string{"m/s", "km/h", "mph"}

var VolumeUnits = []string{
	"cubic-meter",
	"cubic-kilometer",
	"cubic-centimeter",
	"cubic-millimeter",
	"liter",
	"milliliter",
	"gallon",
	"quart",
	"pint",
	"cup",
	"fluid ounce",
}

var ErrUnknownUnit = errors.New("unknown unit")
