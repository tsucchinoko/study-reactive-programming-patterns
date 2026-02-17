package application

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// PlaceOrderCommand represents the intent to create a new order.
type PlaceOrderCommand struct {
	CustomerID   types.CustomerID
	RestaurantID types.RestaurantID
	Items        []PlaceOrderItem
}

// PlaceOrderItem is a DTO for items within a PlaceOrderCommand.
type PlaceOrderItem struct {
	MenuItemID          types.MenuItemID
	Name                string
	Quantity            int
	UnitPrice           types.Money
	SpecialInstructions string
}

// ConfirmOrderCommand represents the intent to confirm an order.
type ConfirmOrderCommand struct {
	OrderID types.OrderID
}

// CancelOrderCommand represents the intent to cancel an order.
type CancelOrderCommand struct {
	OrderID types.OrderID
	Reason  string
}

// TransitionOrderCommand represents the intent to advance an order to the next status.
type TransitionOrderCommand struct {
	OrderID   types.OrderID
	NewStatus string
}
