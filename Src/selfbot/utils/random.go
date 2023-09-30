package utils

import (
	"fmt"
	"math/rand"
)

type RandomUtils struct{}

func (RandomUtils) Number(min, max int) int {
	return rand.Intn(max-min) + min
}

func (RandomUtils) Percentage() string {
	return fmt.Sprintf("%d%%", Random.Number(0, 100))
}

func (RandomUtils) String(length int, alphabet ...string) string {
	var alphabet_ string
	if len(alphabet) > 0 {
		alphabet_ = alphabet[0]
	} else {
		alphabet_ = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet_[Random.Number(0, len(alphabet_))]
	}
	return string(b)
}

// RandomChoice has to be outside the Random struct because of Golang limitations
func RandomChoice[T comparable](choices []T) T {
	return choices[rand.Intn(len(choices))]
}

func RandomShuffle[T comparable](slice []T) []T {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

var Random = &RandomUtils{}
