package types

import "testing"

func TestNewMoney(t *testing.T) {
	m := NewMoney(1000, "USD")
	if m.Amount() != 1000 {
		t.Errorf("Amount() = %d, want 1000", m.Amount())
	}
	if m.Currency() != "USD" {
		t.Errorf("Currency() = %s, want USD", m.Currency())
	}
}

func TestJPY(t *testing.T) {
	m := JPY(500)
	if m.Amount() != 500 {
		t.Errorf("Amount() = %d, want 500", m.Amount())
	}
	if m.Currency() != "JPY" {
		t.Errorf("Currency() = %s, want JPY", m.Currency())
	}
}

func TestMoneyAdd(t *testing.T) {
	a := JPY(300)
	b := JPY(200)
	sum := a.Add(b)
	if sum.Amount() != 500 {
		t.Errorf("Add: Amount() = %d, want 500", sum.Amount())
	}
	if sum.Currency() != "JPY" {
		t.Errorf("Add: Currency() = %s, want JPY", sum.Currency())
	}
	// 元の値が変更されていないことを確認（不変性）
	if a.Amount() != 300 {
		t.Error("Add should not mutate the original")
	}
}

func TestMoneySubtract(t *testing.T) {
	a := JPY(500)
	b := JPY(200)
	diff := a.Subtract(b)
	if diff.Amount() != 300 {
		t.Errorf("Subtract: Amount() = %d, want 300", diff.Amount())
	}
}

func TestMoneyMultiply(t *testing.T) {
	m := JPY(150)
	result := m.Multiply(3)
	if result.Amount() != 450 {
		t.Errorf("Multiply: Amount() = %d, want 450", result.Amount())
	}
	if m.Amount() != 150 {
		t.Error("Multiply should not mutate the original")
	}
}

func TestMoneyIsZero(t *testing.T) {
	if !JPY(0).IsZero() {
		t.Error("JPY(0).IsZero() should be true")
	}
	if JPY(1).IsZero() {
		t.Error("JPY(1).IsZero() should be false")
	}
}

func TestMoneyIsPositive(t *testing.T) {
	if !JPY(100).IsPositive() {
		t.Error("JPY(100).IsPositive() should be true")
	}
	if JPY(0).IsPositive() {
		t.Error("JPY(0).IsPositive() should be false")
	}
	if JPY(-1).IsPositive() {
		t.Error("JPY(-1).IsPositive() should be false")
	}
}

func TestMoneyEqual(t *testing.T) {
	if !JPY(100).Equal(JPY(100)) {
		t.Error("JPY(100) should equal JPY(100)")
	}
	if JPY(100).Equal(JPY(200)) {
		t.Error("JPY(100) should not equal JPY(200)")
	}
	if JPY(100).Equal(NewMoney(100, "USD")) {
		t.Error("JPY(100) should not equal USD(100)")
	}
}

func TestMoneyString(t *testing.T) {
	if got := JPY(1500).String(); got != "¥1500" {
		t.Errorf("String() = %q, want \"¥1500\"", got)
	}
	if got := NewMoney(1000, "USD").String(); got != "USD 1000" {
		t.Errorf("String() = %q, want \"USD 1000\"", got)
	}
}

func TestMoneyAddCurrencyMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add with currency mismatch should panic")
		}
	}()
	JPY(100).Add(NewMoney(100, "USD"))
}

func TestMoneySubtractCurrencyMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Subtract with currency mismatch should panic")
		}
	}()
	JPY(100).Subtract(NewMoney(50, "USD"))
}
