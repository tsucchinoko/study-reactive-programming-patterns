package application

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// PlaceOrderCommand は新しい注文を作成する意図を表す。
type PlaceOrderCommand struct {
	CustomerID   types.CustomerID
	RestaurantID types.RestaurantID
	Items        []PlaceOrderItem
}

// PlaceOrderItem はPlaceOrderCommand内のアイテム用DTO。
type PlaceOrderItem struct {
	MenuItemID          types.MenuItemID
	Name                string
	Quantity            int
	UnitPrice           types.Money
	SpecialInstructions string
}

// ConfirmOrderCommand は注文を確定する意図を表す。
type ConfirmOrderCommand struct {
	OrderID types.OrderID
}

// CancelOrderCommand は注文をキャンセルする意図を表す。
type CancelOrderCommand struct {
	OrderID types.OrderID
	Reason  string
}

// TransitionOrderCommand は注文を次のステータスに進める意図を表す。
type TransitionOrderCommand struct {
	OrderID   types.OrderID
	NewStatus string
}
