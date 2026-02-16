package result

// Result represents either a success value of type T or an error.
// It enforces explicit error handling through Map/FlatMap/Match.
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Unit represents a void success value.
type Unit struct{}

// Ok creates a successful Result.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// FromError converts a (value, error) pair into a Result.
func FromError[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// OkUnit creates a successful Result[Unit].
func OkUnit() Result[Unit] {
	return Ok(Unit{})
}

// ErrUnit creates a failed Result[Unit].
func ErrUnit(err error) Result[Unit] {
	return Err[Unit](err)
}

// IsOk returns true if the Result is a success.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr returns true if the Result is a failure.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// Unwrap returns the success value or panics.
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic("called Unwrap on an Err Result")
	}
	return r.value
}

// UnwrapErr returns the error or panics.
func (r Result[T]) UnwrapErr() error {
	if r.ok {
		panic("called UnwrapErr on an Ok Result")
	}
	return r.err
}

// UnwrapOr returns the success value or the provided default.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.ok {
		return r.value
	}
	return defaultVal
}

// Match calls onOk if success, onErr if failure.
func (r Result[T]) Match(onOk func(T), onErr func(error)) {
	if r.ok {
		onOk(r.value)
	} else {
		onErr(r.err)
	}
}

// Map transforms the success value using f. Errors pass through unchanged.
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if !r.ok {
		return Err[U](r.err)
	}
	return Ok(f(r.value))
}

// FlatMap transforms the success value using f which itself returns a Result.
func FlatMap[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if !r.ok {
		return Err[U](r.err)
	}
	return f(r.value)
}

// MapErr transforms the error value using f. Successes pass through unchanged.
func MapErr[T any](r Result[T], f func(error) error) Result[T] {
	if r.ok {
		return r
	}
	return Err[T](f(r.err))
}
