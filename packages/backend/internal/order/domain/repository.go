package domain

import (
	"context"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// OrderRepository defines the persistence interface for Order aggregates.
// Implementation lives in infrastructure/ — domain knows nothing about storage.
type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	FindByID(ctx context.Context, id types.OrderID) (Order, error)
	FindByCustomerID(ctx context.Context, customerID types.CustomerID) ([]Order, error)
	FindByStatus(ctx context.Context, status OrderStatus) ([]Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]Order, error)
}
