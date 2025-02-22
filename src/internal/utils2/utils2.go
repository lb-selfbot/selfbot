package utils2

import (
	"slices"

	"golang.org/x/exp/constraints"
)

// Clamp constrains a value to be within the specified minimum and maximum bounds.
func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Limit returns a new slice containing at most n elements from the input slice.
// If n is greater than the length of the slice, returns the entire slice.
func Limit[T any](slice []T, n int) []T {
	return slice[:min(len(slice), n)]
}

// Map applies a transformation function to each element in a slice and returns a new slice
func Map[T, R any](slice []T, f func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// MapIndexed applies a transformation function that takes both the index and value
func MapIndexed[T, R any](slice []T, f func(int, T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(i, v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](slice []T, pred func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single value using an accumulator function
func Reduce[T, R any](slice []T, initial R, f func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = f(result, v)
	}
	return result
}

// Any returns true if any element satisfies the predicate
func Any[T any](slice []T, pred func(T) bool) bool {
	return slices.ContainsFunc(slice, pred)
}

// All returns true if all elements satisfy the predicate
func All[T any](slice []T, pred func(T) bool) bool {
	for _, v := range slice {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Contains checks if an element exists in a slice
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// Unique returns a new slice with duplicate elements removed
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
