package fp

// Map applies f to each element in a slice, returning a new slice.
func Map[T any, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}

// Filter returns a new slice containing only elements where predicate returns true.
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce folds a slice into a single value using the accumulator function.
func Reduce[T any, U any](items []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, item := range items {
		acc = f(acc, item)
	}
	return acc
}

// Find returns the first element matching the predicate.
func Find[T any](items []T, predicate func(T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if any element matches the predicate.
func Any[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All returns true if all elements match the predicate.
func All[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Pipe2 composes two functions: Pipe2(f, g)(x) = g(f(x))
func Pipe2[A any, B any, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe3 composes three functions: Pipe3(f, g, h)(x) = h(g(f(x)))
func Pipe3[A any, B any, C any, D any](f func(A) B, g func(B) C, h func(C) D) func(A) D {
	return func(a A) D {
		return h(g(f(a)))
	}
}

// GroupBy groups elements by a key function.
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, item := range items {
		key := keyFn(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// FlatMap applies f to each element and flattens the results.
func FlatMap[T any, U any](items []T, f func(T) []U) []U {
	result := make([]U, 0)
	for _, item := range items {
		result = append(result, f(item)...)
	}
	return result
}
