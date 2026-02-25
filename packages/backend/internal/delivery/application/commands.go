package application

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// AssignDriverCommand はドライバーを注文にアサインする意図を表す。
type AssignDriverCommand struct {
	OrderID        types.OrderID
	PickupLocation domain.Location
}

// UpdateLocationCommand はドライバーの位置を更新する意図を表す。
type UpdateLocationCommand struct {
	DriverID types.DriverID
	Location domain.Location
}

// CompleteDeliveryCommand は配達を完了する意図を表す。
type CompleteDeliveryCommand struct {
	OrderID types.OrderID
}

// PickUpOrderCommand はドライバーが注文をピックアップした意図を表す。
type PickUpOrderCommand struct {
	OrderID types.OrderID
}
