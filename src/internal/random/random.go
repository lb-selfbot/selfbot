package random

import (
	"math/rand"
	"strconv"
)

const (
	AlphabetLowercase = "abcdefghijklmnopqrstuvwxyz"
	AlphabetUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetNumbers   = "0123456789"
	AlphabetAll       = AlphabetLowercase + AlphabetUppercase + AlphabetNumbers
)

// Int returns a random integer between min and max.
func Int(min, max int) int {
	return rand.Intn(max-min) + min
}

// Percentage returns a random percentage.
func Percentage() string {
	return strconv.Itoa(Int(0, 100)) + "%"
}

// String returns a random string of length n.
// If charsets is empty, it uses AlphabetAll.
func String(n int, charsets ...string) string {
	if len(charsets) == 0 {
		charsets = append(charsets, AlphabetAll)
	}

	charset := charsets[0]

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Choice returns a random element from the choices slice.
func Choice[T comparable](choices []T) T {
	return choices[rand.Intn(len(choices))]
}

// Shuffle returns a shuffled slice.
func Shuffle[T comparable](slice []T) []T {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}
