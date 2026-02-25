package fp

// Map はスライスの各要素に f を適用し、新しいスライスを返す。
func Map[T any, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}

// Filter は predicate が true を返す要素のみを含む新しいスライスを返す。
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce はアキュムレータ関数を使ってスライスを単一の値に畳み込む。
func Reduce[T any, U any](items []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, item := range items {
		acc = f(acc, item)
	}
	return acc
}

// Find は predicate に一致する最初の要素を返す。
func Find[T any](items []T, predicate func(T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Any はいずれかの要素が predicate に一致する場合に true を返す。
func Any[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All はすべての要素が predicate に一致する場合に true を返す。
func All[T any](items []T, predicate func(T) bool) bool {
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Pipe2 は2つの関数を合成する: Pipe2(f, g)(x) = g(f(x))
func Pipe2[A any, B any, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe3 は3つの関数を合成する: Pipe3(f, g, h)(x) = h(g(f(x)))
func Pipe3[A any, B any, C any, D any](f func(A) B, g func(B) C, h func(C) D) func(A) D {
	return func(a A) D {
		return h(g(f(a)))
	}
}

// GroupBy はキー関数を使って要素をグループ化する。
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, item := range items {
		key := keyFn(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// FlatMap は各要素に f を適用し、結果を平坦化して返す。
func FlatMap[T any, U any](items []T, f func(T) []U) []U {
	result := make([]U, 0)
	for _, item := range items {
		result = append(result, f(item)...)
	}
	return result
}
