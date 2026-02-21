package domain

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// OrderItem は注文内の1つの明細を表す値オブジェクト。
// イミュータブル — 全フィールドは非公開でアクセサメソッドを持つ。
type OrderItem struct {
	menuItemID          types.MenuItemID
	name                string
	quantity            int
	unitPrice           types.Money
	specialInstructions string
}

// NewOrderItem は新しいOrderItem値オブジェクトを作成する。
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

// Subtotal はunitPrice * quantityを返す。純粋関数。
func (i OrderItem) Subtotal() types.Money {
	return i.unitPrice.Multiply(i.quantity)
}
