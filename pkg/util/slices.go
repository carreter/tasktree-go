// Package util provides various helper functions.
//
// I know that Map/Reduce/Filter isn't idiomatic Go, but I still find them very useful!
package util

// Map creates a new slice from the result of calling a function on each element
// of an input slice.
func Map[T any, R any](l []T, f func(T) R) []R {
	res := make([]R, len(l))
	for i, elem := range l {
		res[i] = f(elem)
	}
	return res
}

// ForEach applies a function to each element in a slice.
func ForEach[T any](l []T, f func(T)) {
	for _, elem := range l {
		f(elem)
	}
}

// Filter iterates over a slice and returns a slice containing only the elements
// for which the predicate is true.
func Filter[T any](l []T, f func(T) bool) []T {
	res := make([]T, 0)
	for _, elem := range l {
		if f(elem) {
			res = append(res, elem)
		}
	}
	return res
}

// Remove removes all instances of an element from a slice.
func Remove[T comparable](l []T, toRemove T) []T {
	return Filter(l, func(elem T) bool {
		return elem != toRemove
	})
}

// Reduce applies an accumulator function from left to right on a slice with a given initial value.
func Reduce[T any, R any](l []T, f func(R, T) R, acc R) R {
	for _, elem := range l {
		acc = f(acc, elem)
	}
	return acc
}

// Contains checks if a slice contains a given element.
func Contains[T comparable](l []T, query T) bool {
	for _, elem := range l {
		if elem == query {
			return true
		}
	}
	return false
}
