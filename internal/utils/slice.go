package utils

type predicate[T any] func(*T) bool

func Find[T any](slice []T, pred predicate[T]) (*T, bool) {
	for i := range slice {
		if pred(&slice[i]) {
			return &slice[i], true
		}
	}

	return nil, false
}

func Filter[T any](slice []T, pred predicate[T]) []T {
	subset := []T{}

	for i := range slice {
		if pred(&slice[i]) {
			subset = append(subset, slice[i])
		}
	}

	return subset
}

func ForEach[T any](slice []T, f func(T)) {
	for _, value := range slice {
		f(value)
	}
}
