package formatting

import (
	"strconv"
	"strings"
	"time"
)

// TimestampFlag is a flag for formatting timestamps.
type TimestampFlag string

const (
	// TimestampFlagShortTime is a flag for short time format.
	TimestampFlagShortTime TimestampFlag = "t"

	// TimestampFlagLongTime is a flag for long time format.
	TimestampFlagLongTime TimestampFlag = "T"

	// TimestampFlagShortDate is a flag for short date format.
	TimestampFlagShortDate TimestampFlag = "d"

	// TimestampFlagLongDate is a flag for long date format.
	TimestampFlagLongDate TimestampFlag = "D"

	// TimestampFlagShortDateTime is a flag for short date and time format.
	TimestampFlagShortDateTime TimestampFlag = "f"

	// TimestampFlagLongDateTime is a flag for long date and time format.
	TimestampFlagLongDateTime TimestampFlag = "F"

	// TimestampFlagRelative is a flag for relative time format.
	TimestampFlagRelative TimestampFlag = "R"
)

// Timestamp formats a time.Time value into a Discord timestamp string.
// Only the first flag is used if multiple are provided.
// Format: <t:unix> without flags; <t:unix:flag> with a flag.
func Timestamp(t time.Time, flags ...TimestampFlag) string {
	var builder strings.Builder

	builder.WriteString("<t:")
	builder.WriteString(strconv.FormatInt(t.Unix(), 10))
	if len(flags) == 1 {
		builder.WriteString(":")
		builder.WriteString(string(flags[0]))
	}
	builder.WriteString(">")

	return builder.String()
}
