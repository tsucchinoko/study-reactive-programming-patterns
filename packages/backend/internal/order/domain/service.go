package domain

import (
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// --- Pure domain service functions ---
// These functions contain business logic with no side effects.
// They are easily testable without mocks.

// ValidateOrderItems checks that all items have valid quantities and prices.
func ValidateOrderItems(items []OrderItem) result.Result[[]OrderItem] {
	for _, item := range items {
		if item.Quantity() <= 0 {
			return result.Err[[]OrderItem](
				fmt.Errorf("item %s: quantity must be positive", item.Name()),
			)
		}
		if !item.UnitPrice().IsPositive() {
			return result.Err[[]OrderItem](
				fmt.Errorf("item %s: price must be positive", item.Name()),
			)
		}
	}
	return result.Ok(items)
}

// BuildOrder validates items and creates a new Order. Pure function.
func BuildOrder(
	id types.OrderID,
	customerID types.CustomerID,
	restaurantID types.RestaurantID,
	items []OrderItem,
	placedAt types.Timestamp,
) result.Result[Order] {
	return result.FlatMap(
		ValidateOrderItems(items),
		func(validItems []OrderItem) result.Result[Order] {
			return NewOrder(id, customerID, restaurantID, validItems, placedAt)
		},
	)
}
