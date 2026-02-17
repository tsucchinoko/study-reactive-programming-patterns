package types

import "fmt"

// Money represents a monetary amount in the smallest unit (cents) to avoid floating-point issues.
// Money is a value object — all operations return new instances.
type Money struct {
	amount   int64
	currency string
}

// NewMoney creates a Money value object.
func NewMoney(amount int64, currency string) Money {
	return Money{amount: amount, currency: currency}
}

// JPY creates a Money value in Japanese Yen.
func JPY(amount int64) Money {
	return Money{amount: amount, currency: "JPY"}
}

// Amount returns the raw amount in smallest unit.
func (m Money) Amount() int64 { return m.amount }

// Currency returns the currency code.
func (m Money) Currency() string { return m.currency }

// Add returns a new Money that is the sum of m and other.
// Panics if currencies don't match.
func (m Money) Add(other Money) Money {
	m.assertSameCurrency(other)
	return Money{amount: m.amount + other.amount, currency: m.currency}
}

// Subtract returns a new Money that is m minus other.
func (m Money) Subtract(other Money) Money {
	m.assertSameCurrency(other)
	return Money{amount: m.amount - other.amount, currency: m.currency}
}

// Multiply returns a new Money multiplied by the quantity.
func (m Money) Multiply(quantity int) Money {
	return Money{amount: m.amount * int64(quantity), currency: m.currency}
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive returns true if the amount is greater than zero.
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// Equal returns true if both Money values are identical.
func (m Money) Equal(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns a human-readable representation.
func (m Money) String() string {
	if m.currency == "JPY" {
		return fmt.Sprintf("¥%d", m.amount)
	}
	return fmt.Sprintf("%s %d", m.currency, m.amount)
}

func (m Money) assertSameCurrency(other Money) {
	if m.currency != other.currency {
		panic(fmt.Sprintf("currency mismatch: %s vs %s", m.currency, other.currency))
	}
}
