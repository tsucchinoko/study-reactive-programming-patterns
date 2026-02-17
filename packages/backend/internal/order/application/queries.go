package application

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// GetOrderQuery retrieves a single order by ID.
type GetOrderQuery struct {
	OrderID types.OrderID
}

// ListOrdersQuery retrieves orders with pagination.
type ListOrdersQuery struct {
	Limit  int
	Offset int
}

// ListOrdersByStatusQuery retrieves orders filtered by status.
type ListOrdersByStatusQuery struct {
	Status string
}

// ListOrdersByCustomerQuery retrieves orders for a specific customer.
type ListOrdersByCustomerQuery struct {
	CustomerID types.CustomerID
}
