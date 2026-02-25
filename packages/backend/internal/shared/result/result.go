package result

// Result は型 T の成功値またはエラーのいずれかを表す。
// Map/FlatMap/Match を通じて明示的なエラーハンドリングを強制する。
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Unit は void の成功値を表す。
type Unit struct{}

// Ok は成功した Result を生成する。
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

// Err は失敗した Result を生成する。
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// FromError は (value, error) のペアを Result に変換する。
func FromError[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// OkUnit は成功した Result[Unit] を生成する。
func OkUnit() Result[Unit] {
	return Ok(Unit{})
}

// ErrUnit は失敗した Result[Unit] を生成する。
func ErrUnit(err error) Result[Unit] {
	return Err[Unit](err)
}

// IsOk は Result が成功の場合に true を返す。
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr は Result が失敗の場合に true を返す。
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// Unwrap は成功値を返す。失敗の場合はパニックする。
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic("called Unwrap on an Err Result")
	}
	return r.value
}

// UnwrapErr はエラーを返す。成功の場合はパニックする。
func (r Result[T]) UnwrapErr() error {
	if r.ok {
		panic("called UnwrapErr on an Ok Result")
	}
	return r.err
}

// UnwrapOr は成功値を返す。失敗の場合は指定したデフォルト値を返す。
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.ok {
		return r.value
	}
	return defaultVal
}

// Match は成功の場合に onOk を、失敗の場合に onErr を呼び出す。
func (r Result[T]) Match(onOk func(T), onErr func(error)) {
	if r.ok {
		onOk(r.value)
	} else {
		onErr(r.err)
	}
}

// Map は f を使って成功値を変換する。エラーはそのまま通過する。
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if !r.ok {
		return Err[U](r.err)
	}
	return Ok(f(r.value))
}

// FlatMap は Result を返す f を使って成功値を変換する。
func FlatMap[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if !r.ok {
		return Err[U](r.err)
	}
	return f(r.value)
}

// MapErr は f を使ってエラー値を変換する。成功はそのまま通過する。
func MapErr[T any](r Result[T], f func(error) error) Result[T] {
	if r.ok {
		return r
	}
	return Err[T](f(r.err))
}
