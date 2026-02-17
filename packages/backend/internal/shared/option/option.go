package option

// Option represents an optional value of type T.
// It replaces nil pointers with explicit presence/absence semantics.
type Option[T any] struct {
	value *T
}

// Some creates an Option containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

// None creates an empty Option.
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// FromPtr creates an Option from a pointer.
func FromPtr[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone returns true if the Option is empty.
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Unwrap returns the value or panics if empty.
func (o Option[T]) Unwrap() T {
	if o.value == nil {
		panic("called Unwrap on a None Option")
	}
	return *o.value
}

// UnwrapOr returns the value or the provided default.
func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.value == nil {
		return defaultVal
	}
	return *o.value
}

// ToPtr returns the underlying pointer (nil if None).
func (o Option[T]) ToPtr() *T {
	return o.value
}

// Match calls onSome if present, onNone if empty.
func (o Option[T]) Match(onSome func(T), onNone func()) {
	if o.value != nil {
		onSome(*o.value)
	} else {
		onNone()
	}
}

// Map transforms the contained value using f. None passes through unchanged.
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return Some(f(*o.value))
}

// FlatMap transforms the contained value using f which itself returns an Option.
func FlatMap[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return f(*o.value)
}

// Filter returns None if the predicate returns false.
func Filter[T any](o Option[T], predicate func(T) bool) Option[T] {
	if o.value != nil && predicate(*o.value) {
		return o
	}
	return None[T]()
}
