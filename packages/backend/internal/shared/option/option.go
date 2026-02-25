package option

// Option は型 T のオプショナルな値を表す。
// nil ポインタを明示的な有無のセマンティクスに置き換える。
type Option[T any] struct {
	value *T
}

// Some は値を含む Option を生成する。
func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

// None は空の Option を生成する。
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// FromPtr はポインタから Option を生成する。
func FromPtr[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

// IsSome は Option が値を含む場合に true を返す。
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone は Option が空の場合に true を返す。
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Unwrap は値を返す。Option が空の場合はパニックする。
func (o Option[T]) Unwrap() T {
	if o.value == nil {
		panic("called Unwrap on a None Option")
	}
	return *o.value
}

// UnwrapOr は値を返す。Option が空の場合は指定したデフォルト値を返す。
func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.value == nil {
		return defaultVal
	}
	return *o.value
}

// ToPtr は内部ポインタを返す（None の場合は nil）。
func (o Option[T]) ToPtr() *T {
	return o.value
}

// Match は値がある場合に onSome を、空の場合に onNone を呼び出す。
func (o Option[T]) Match(onSome func(T), onNone func()) {
	if o.value != nil {
		onSome(*o.value)
	} else {
		onNone()
	}
}

// Map は f を使って内包する値を変換する。None はそのまま通過する。
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return Some(f(*o.value))
}

// FlatMap は Option を返す f を使って内包する値を変換する。
func FlatMap[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return f(*o.value)
}

// Filter は predicate が false を返す場合に None を返す。
func Filter[T any](o Option[T], predicate func(T) bool) Option[T] {
	if o.value != nil && predicate(*o.value) {
		return o
	}
	return None[T]()
}
