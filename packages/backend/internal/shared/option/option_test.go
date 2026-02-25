package option

import "testing"

func TestSome(t *testing.T) {
	o := Some(42)
	if !o.IsSome() {
		t.Error("Some should be IsSome")
	}
	if o.IsNone() {
		t.Error("Some should not be IsNone")
	}
	if o.Unwrap() != 42 {
		t.Errorf("Unwrap() = %d, want 42", o.Unwrap())
	}
}

func TestNone(t *testing.T) {
	o := None[int]()
	if o.IsSome() {
		t.Error("None should not be IsSome")
	}
	if !o.IsNone() {
		t.Error("None should be IsNone")
	}
}

func TestFromPtr_NonNil(t *testing.T) {
	v := 42
	o := FromPtr(&v)
	if !o.IsSome() {
		t.Error("FromPtr with non-nil should be Some")
	}
	if o.Unwrap() != 42 {
		t.Errorf("Unwrap() = %d, want 42", o.Unwrap())
	}
}

func TestFromPtr_Nil(t *testing.T) {
	o := FromPtr[int](nil)
	if !o.IsNone() {
		t.Error("FromPtr with nil should be None")
	}
}

func TestUnwrapPanicsOnNone(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Unwrap on None should panic")
		}
	}()
	None[int]().Unwrap()
}

func TestUnwrapOr_Some(t *testing.T) {
	o := Some(42)
	if o.UnwrapOr(0) != 42 {
		t.Errorf("UnwrapOr on Some should return value, got %d", o.UnwrapOr(0))
	}
}

func TestUnwrapOr_None(t *testing.T) {
	o := None[int]()
	if o.UnwrapOr(99) != 99 {
		t.Errorf("UnwrapOr on None should return default, got %d", o.UnwrapOr(99))
	}
}

func TestToPtr_Some(t *testing.T) {
	o := Some(42)
	p := o.ToPtr()
	if p == nil {
		t.Fatal("ToPtr on Some should return non-nil pointer")
	}
	if *p != 42 {
		t.Errorf("*ToPtr() = %d, want 42", *p)
	}
}

func TestToPtr_None(t *testing.T) {
	o := None[int]()
	if o.ToPtr() != nil {
		t.Error("ToPtr on None should return nil")
	}
}

func TestMatch_Some(t *testing.T) {
	var called bool
	Some(42).Match(
		func(v int) { called = true },
		func() { t.Error("onNone should not be called for Some") },
	)
	if !called {
		t.Error("onSome should be called")
	}
}

func TestMatch_None(t *testing.T) {
	var called bool
	None[int]().Match(
		func(v int) { t.Error("onSome should not be called for None") },
		func() { called = true },
	)
	if !called {
		t.Error("onNone should be called")
	}
}

func TestMap_Some(t *testing.T) {
	o := Map(Some(10), func(v int) string { return "ten" })
	if !o.IsSome() {
		t.Error("Map on Some should return Some")
	}
	if o.Unwrap() != "ten" {
		t.Errorf("Map result = %q, want \"ten\"", o.Unwrap())
	}
}

func TestMap_None(t *testing.T) {
	o := Map(None[int](), func(v int) string { return "unreachable" })
	if !o.IsNone() {
		t.Error("Map on None should return None")
	}
}

func TestFlatMap_Some(t *testing.T) {
	o := FlatMap(Some(10), func(v int) Option[string] {
		return Some("ok")
	})
	if o.Unwrap() != "ok" {
		t.Errorf("FlatMap result = %q, want \"ok\"", o.Unwrap())
	}
}

func TestFlatMap_SomeToNone(t *testing.T) {
	o := FlatMap(Some(10), func(v int) Option[string] {
		return None[string]()
	})
	if !o.IsNone() {
		t.Error("FlatMap Some -> None should return None")
	}
}

func TestFlatMap_None(t *testing.T) {
	o := FlatMap(None[int](), func(v int) Option[string] {
		t.Error("function should not be called on None")
		return Some("unreachable")
	})
	if !o.IsNone() {
		t.Error("FlatMap on None should return None")
	}
}

func TestFilter_SomeMatch(t *testing.T) {
	o := Filter(Some(10), func(v int) bool { return v > 5 })
	if !o.IsSome() {
		t.Error("Filter with matching predicate should return Some")
	}
	if o.Unwrap() != 10 {
		t.Errorf("Filter result = %d, want 10", o.Unwrap())
	}
}

func TestFilter_SomeNoMatch(t *testing.T) {
	o := Filter(Some(3), func(v int) bool { return v > 5 })
	if !o.IsNone() {
		t.Error("Filter with non-matching predicate should return None")
	}
}

func TestFilter_None(t *testing.T) {
	o := Filter(None[int](), func(v int) bool { return true })
	if !o.IsNone() {
		t.Error("Filter on None should return None")
	}
}
