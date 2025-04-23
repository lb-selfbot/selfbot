package unsafeutil

import "unsafe"

// StringToByteSlice converts a string to a byte slice without copying the data.
func StringToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ByteSliceToString converts a byte slice to a string without copying the data.
func ByteSliceToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
