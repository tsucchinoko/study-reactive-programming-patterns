package domain

import (
	"context"

	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// DriverRepository はドライバーの永続化インターフェース。
type DriverRepository interface {
	Save(ctx context.Context, driver Driver) error
	FindByID(ctx context.Context, id types.DriverID) (Driver, error)
	FindByStatus(ctx context.Context, status DriverStatus) ([]Driver, error)
	FindAll(ctx context.Context) ([]Driver, error)
}

// AssignmentRepository は配達アサインメントの永続化インターフェース。
type AssignmentRepository interface {
	Save(ctx context.Context, assignment DeliveryAssignment) error
	FindByID(ctx context.Context, id types.AssignmentID) (DeliveryAssignment, error)
	FindByOrderID(ctx context.Context, orderID types.OrderID) (DeliveryAssignment, error)
	FindByDriverID(ctx context.Context, driverID types.DriverID) ([]DeliveryAssignment, error)
}
