package domain

import (
	"fmt"

	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/result"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// --- 純粋なドメインサービス関数 ---
// これらの関数は副作用のないビジネスロジックを含む。
// モックなしで簡単にテスト可能。

// ValidateOrderItems は全アイテムの数量と価格が有効かを検証する。
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

// BuildOrder はアイテムを検証し新しいOrderを作成する。純粋関数。
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
