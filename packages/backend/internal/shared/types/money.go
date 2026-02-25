package types

import "fmt"

// Money は浮動小数点の問題を避けるため、最小単位（銭）で金額を表す値オブジェクト。
// Money は値オブジェクトであり、すべての操作は新しいインスタンスを返す。
type Money struct {
	amount   int64
	currency string
}

// NewMoney は Money 値オブジェクトを生成する。
func NewMoney(amount int64, currency string) Money {
	return Money{amount: amount, currency: currency}
}

// JPY は日本円の Money 値を生成する。
func JPY(amount int64) Money {
	return Money{amount: amount, currency: "JPY"}
}

// Amount は最小単位での生の金額を返す。
func (m Money) Amount() int64 { return m.amount }

// Currency は通貨コードを返す。
func (m Money) Currency() string { return m.currency }

// Add は m と other の合計を表す新しい Money を返す。
// 通貨が一致しない場合はパニックする。
func (m Money) Add(other Money) Money {
	m.assertSameCurrency(other)
	return Money{amount: m.amount + other.amount, currency: m.currency}
}

// Subtract は m から other を引いた新しい Money を返す。
func (m Money) Subtract(other Money) Money {
	m.assertSameCurrency(other)
	return Money{amount: m.amount - other.amount, currency: m.currency}
}

// Multiply は数量を乗算した新しい Money を返す。
func (m Money) Multiply(quantity int) Money {
	return Money{amount: m.amount * int64(quantity), currency: m.currency}
}

// IsZero は金額がゼロの場合に true を返す。
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive は金額がゼロより大きい場合に true を返す。
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// Equal は両方の Money 値が等しい場合に true を返す。
func (m Money) Equal(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String は人間が読みやすい表現を返す。
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
