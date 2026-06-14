package list

type predicate[T any] func(*T) bool

func Find[T any](list []T, pred predicate[T]) (*T, bool) {
	for i := range list {
		if pred(&list[i]) {
			return &list[i], true
		}
	}

	return nil, false
}

func Filter[T any](list []T, pred predicate[T]) []T {
	subset := []T{}

	for i := range list {
		if pred(&list[i]) {
			subset = append(subset, list[i])
		}
	}

	return subset
}
