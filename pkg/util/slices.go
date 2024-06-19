package util

// I know that Map/Reduce/Filter isn't idiomatic Go, but I still find them very useful!

func Map[T any, R any](l []T, f func(T) R) []R {
	res := make([]R, len(l))
	for i, elem := range l {
		res[i] = f(elem)
	}
	return res
}

func ForEach[T any](l []T, f func(T)) {
	for _, elem := range l {
		f(elem)
	}
}

func Filter[T any](l []T, f func(T) bool) []T {
	res := make([]T, 0)
	for _, elem := range l {
		if f(elem) {
			res = append(res, elem)
		}
	}
	return res
}

func Remove[T comparable](l []T, toRemove T) []T {
	return Filter(l, func(elem T) bool {
		return elem != toRemove
	})
}

func Reduce[T any, R any](l []T, f func(R, T) R, acc R) R {
	for _, elem := range l {
		acc = f(acc, elem)
	}
	return acc
}

func Contains[T comparable](l []T, query T) bool {
	for _, elem := range l {
		if elem == query {
			return true
		}
	}
	return false
}
