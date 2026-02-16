package domain

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// OrderItem is a value object representing a single line item in an order.
// Immutable — all fields are unexported with accessor methods.
type OrderItem struct {
	menuItemID          types.MenuItemID
	name                string
	quantity            int
	unitPrice           types.Money
	specialInstructions string
}

// NewOrderItem creates a new OrderItem value object.
func NewOrderItem(menuItemID types.MenuItemID, name string, quantity int, unitPrice types.Money, instructions string) OrderItem {
	return OrderItem{
		menuItemID:          menuItemID,
		name:                name,
		quantity:            quantity,
		unitPrice:           unitPrice,
		specialInstructions: instructions,
	}
}

func (i OrderItem) MenuItemID() types.MenuItemID { return i.menuItemID }
func (i OrderItem) Name() string                 { return i.name }
func (i OrderItem) Quantity() int                { return i.quantity }
func (i OrderItem) UnitPrice() types.Money       { return i.unitPrice }
func (i OrderItem) SpecialInstructions() string  { return i.specialInstructions }

// Subtotal returns unitPrice * quantity. Pure function.
func (i OrderItem) Subtotal() types.Money {
	return i.unitPrice.Multiply(i.quantity)
}
