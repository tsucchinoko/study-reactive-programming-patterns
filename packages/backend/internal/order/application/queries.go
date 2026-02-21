package application

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// GetOrderQuery はIDで単一の注文を取得する。
type GetOrderQuery struct {
	OrderID types.OrderID
}

// ListOrdersQuery はページネーション付きで注文を取得する。
type ListOrdersQuery struct {
	Limit  int
	Offset int
}

// ListOrdersByStatusQuery はステータスでフィルタリングした注文を取得する。
type ListOrdersByStatusQuery struct {
	Status string
}

// ListOrdersByCustomerQuery は特定の顧客の注文を取得する。
type ListOrdersByCustomerQuery struct {
	CustomerID types.CustomerID
}
