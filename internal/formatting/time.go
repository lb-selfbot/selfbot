package formatting

import (
	"fmt"
	"strings"
	"time"
)

const (
	secondsPerMinute = 60
	secondsPerHour   = 3600  // 60 * 60
	secondsPerDay    = 86400 // 3600 * 24
)

// plural formats a unit string based on the value.
func plural(val int64, unit string) string {
	if val == 1 {
		return fmt.Sprintf("%d %s", val, unit)
	}
	return fmt.Sprintf("%d %ss", val, unit)
}

type durationUnits struct {
	day    string
	hour   string
	minute string
	second string
}

// duration formats a time.duration into a human-readable string like "X days Y hours Z minutes W seconds".
func duration(d time.Duration, units durationUnits) string {
	totalSeconds := int64(d.Seconds())

	if totalSeconds <= 0 {
		return "0 seconds"
	}

	days := totalSeconds / secondsPerDay
	remainingSeconds := totalSeconds % secondsPerDay

	hours := remainingSeconds / secondsPerHour
	remainingSeconds %= secondsPerHour

	minutes := remainingSeconds / secondsPerMinute
	seconds := remainingSeconds % secondsPerMinute

	parts := []string{}

	if days > 0 {
		parts = append(parts, plural(days, units.day))
	}
	if hours > 0 {
		parts = append(parts, plural(hours, units.hour))
	}
	if minutes > 0 {
		parts = append(parts, plural(minutes, units.minute))
	}
	if seconds > 0 || len(parts) == 0 { // Always show seconds if it's the only unit or non-zero
		parts = append(parts, plural(seconds, units.second))
	}

	return strings.Join(parts, " ")
}

// DurationLong formats a time.Duration into a human-readable string like "X days Y hours Z minutes W seconds".
func DurationLong(d time.Duration) string {
	return duration(d, durationUnits{
		day:    "day",
		hour:   "hour",
		minute: "minute",
		second: "second",
	})
}

// DurationShort formats a time.Duration into a human-readable string like "Xh Ym Zs".
func DurationShort(d time.Duration) string {
	return duration(d, durationUnits{
		day:    "d",
		hour:   "h",
		minute: "m",
		second: "s",
	})
}
