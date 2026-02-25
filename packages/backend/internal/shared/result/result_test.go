package result

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	r := Ok(42)
	if !r.IsOk() {
		t.Error("Ok should be IsOk")
	}
	if r.IsErr() {
		t.Error("Ok should not be IsErr")
	}
	if r.Unwrap() != 42 {
		t.Errorf("Unwrap() = %d, want 42", r.Unwrap())
	}
}

func TestErr(t *testing.T) {
	e := errors.New("something failed")
	r := Err[int](e)
	if r.IsOk() {
		t.Error("Err should not be IsOk")
	}
	if !r.IsErr() {
		t.Error("Err should be IsErr")
	}
	if r.UnwrapErr() != e {
		t.Errorf("UnwrapErr() = %v, want %v", r.UnwrapErr(), e)
	}
}

func TestFromError_Success(t *testing.T) {
	r := FromError(42, nil)
	if !r.IsOk() {
		t.Error("FromError with nil error should be Ok")
	}
	if r.Unwrap() != 42 {
		t.Errorf("Unwrap() = %d, want 42", r.Unwrap())
	}
}

func TestFromError_Failure(t *testing.T) {
	e := errors.New("fail")
	r := FromError(0, e)
	if !r.IsErr() {
		t.Error("FromError with error should be Err")
	}
	if r.UnwrapErr() != e {
		t.Errorf("UnwrapErr() = %v, want %v", r.UnwrapErr(), e)
	}
}

func TestOkUnit(t *testing.T) {
	r := OkUnit()
	if !r.IsOk() {
		t.Error("OkUnit should be Ok")
	}
}

func TestErrUnit(t *testing.T) {
	e := errors.New("unit error")
	r := ErrUnit(e)
	if !r.IsErr() {
		t.Error("ErrUnit should be Err")
	}
	if r.UnwrapErr() != e {
		t.Errorf("UnwrapErr() = %v, want %v", r.UnwrapErr(), e)
	}
}

func TestUnwrapPanicsOnErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Unwrap on Err should panic")
		}
	}()
	Err[int](errors.New("fail")).Unwrap()
}

func TestUnwrapErrPanicsOnOk(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("UnwrapErr on Ok should panic")
		}
	}()
	Ok(42).UnwrapErr()
}

func TestUnwrapOr_Ok(t *testing.T) {
	r := Ok(42)
	if r.UnwrapOr(0) != 42 {
		t.Errorf("UnwrapOr on Ok should return value, got %d", r.UnwrapOr(0))
	}
}

func TestUnwrapOr_Err(t *testing.T) {
	r := Err[int](errors.New("fail"))
	if r.UnwrapOr(99) != 99 {
		t.Errorf("UnwrapOr on Err should return default, got %d", r.UnwrapOr(99))
	}
}

func TestMatch_Ok(t *testing.T) {
	var called bool
	Ok(42).Match(
		func(v int) { called = true },
		func(err error) { t.Error("onErr should not be called for Ok") },
	)
	if !called {
		t.Error("onOk should be called")
	}
}

func TestMatch_Err(t *testing.T) {
	var called bool
	Err[int](errors.New("fail")).Match(
		func(v int) { t.Error("onOk should not be called for Err") },
		func(err error) { called = true },
	)
	if !called {
		t.Error("onErr should be called")
	}
}

func TestMap_Ok(t *testing.T) {
	r := Map(Ok(10), func(v int) string { return "val" })
	if !r.IsOk() {
		t.Error("Map on Ok should return Ok")
	}
	if r.Unwrap() != "val" {
		t.Errorf("Map result = %q, want \"val\"", r.Unwrap())
	}
}

func TestMap_Err(t *testing.T) {
	e := errors.New("fail")
	r := Map(Err[int](e), func(v int) string { return "val" })
	if !r.IsErr() {
		t.Error("Map on Err should return Err")
	}
	if r.UnwrapErr() != e {
		t.Error("Map on Err should preserve original error")
	}
}

func TestFlatMap_Ok(t *testing.T) {
	r := FlatMap(Ok(10), func(v int) Result[string] {
		return Ok("ok")
	})
	if r.Unwrap() != "ok" {
		t.Errorf("FlatMap result = %q, want \"ok\"", r.Unwrap())
	}
}

func TestFlatMap_OkToErr(t *testing.T) {
	e := errors.New("inner fail")
	r := FlatMap(Ok(10), func(v int) Result[string] {
		return Err[string](e)
	})
	if !r.IsErr() {
		t.Error("FlatMap Ok -> Err should return Err")
	}
}

func TestFlatMap_Err(t *testing.T) {
	e := errors.New("outer fail")
	r := FlatMap(Err[int](e), func(v int) Result[string] {
		t.Error("function should not be called on Err")
		return Ok("unreachable")
	})
	if !r.IsErr() {
		t.Error("FlatMap on Err should return Err")
	}
	if r.UnwrapErr() != e {
		t.Error("FlatMap on Err should preserve original error")
	}
}

func TestMapErr_Ok(t *testing.T) {
	r := MapErr(Ok(42), func(err error) error {
		return errors.New("wrapped")
	})
	if !r.IsOk() {
		t.Error("MapErr on Ok should return Ok")
	}
	if r.Unwrap() != 42 {
		t.Errorf("MapErr on Ok: Unwrap() = %d, want 42", r.Unwrap())
	}
}

func TestMapErr_Err(t *testing.T) {
	original := errors.New("original")
	r := MapErr(Err[int](original), func(err error) error {
		return errors.New("wrapped: " + err.Error())
	})
	if !r.IsErr() {
		t.Error("MapErr on Err should return Err")
	}
	if r.UnwrapErr().Error() != "wrapped: original" {
		t.Errorf("MapErr error = %q, want \"wrapped: original\"", r.UnwrapErr().Error())
	}
}
