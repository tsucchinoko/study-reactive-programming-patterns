package domain

import (
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
)

// OrderStatus は注文のライフサイクルにおける現在の状態を表す。
type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"
	StatusConfirmed  OrderStatus = "CONFIRMED"
	StatusPreparing  OrderStatus = "PREPARING"
	StatusReady      OrderStatus = "READY"
	StatusPickedUp   OrderStatus = "PICKED_UP"
	StatusDelivering OrderStatus = "DELIVERING"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

// validTransitions は許可されたステートマシンの遷移を定義する。
// 純粋なデータ宣言 — 副作用なし。
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusCreated:    {StatusConfirmed, StatusCancelled},
	StatusConfirmed:  {StatusPreparing, StatusCancelled},
	StatusPreparing:  {StatusReady, StatusCancelled},
	StatusReady:      {StatusPickedUp, StatusCancelled},
	StatusPickedUp:   {StatusDelivering},
	StatusDelivering: {StatusDelivered},
	StatusDelivered:  {},
	StatusCancelled:  {},
}

// CanTransitionTo はステータス遷移が有効かどうかを検証する。
// 純粋関数 — 副作用なし。
func CanTransitionTo(current, next OrderStatus) result.Result[result.Unit] {
	allowed, exists := validTransitions[current]
	if !exists {
		return result.Err[result.Unit](fmt.Errorf("unknown order status: %s", current))
	}
	for _, s := range allowed {
		if s == next {
			return result.OkUnit()
		}
	}
	return result.Err[result.Unit](
		fmt.Errorf("invalid transition: %s → %s", current, next),
	)
}

// NextStatuses は現在のステータスから有効な次のステータスのリストを返す。
// 純粋関数。
func NextStatuses(current OrderStatus) []OrderStatus {
	return validTransitions[current]
}

// IsTerminal はこれ以上の遷移が不可能な場合にtrueを返す。
// 純粋関数。
func IsTerminal(status OrderStatus) bool {
	return len(validTransitions[status]) == 0
}

// AllStatuses は全ての注文ステータスを返す。
func AllStatuses() []OrderStatus {
	return []OrderStatus{
		StatusCreated, StatusConfirmed, StatusPreparing, StatusReady,
		StatusPickedUp, StatusDelivering, StatusDelivered, StatusCancelled,
	}
}

func (s OrderStatus) String() string { return string(s) }
