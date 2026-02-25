package domain

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

// --- 配達ドメインイベント ---

// DriverAssigned はドライバーが注文にアサインされた時に発行される。
type DriverAssigned struct {
	events.BaseEvent
	OrderID  string `json:"order_id"`
	DriverID string `json:"driver_id"`
}

func NewDriverAssigned(orderID, driverID string) DriverAssigned {
	return DriverAssigned{
		BaseEvent: events.NewBaseEvent(orderID, "DriverAssigned", "delivery.driver_assigned"),
		OrderID:   orderID,
		DriverID:  driverID,
	}
}

// DriverLocationUpdated はドライバーの位置が更新された時に発行される。
type DriverLocationUpdated struct {
	events.BaseEvent
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewDriverLocationUpdated(driverID string, lat, lng float64) DriverLocationUpdated {
	return DriverLocationUpdated{
		BaseEvent: events.NewBaseEvent(driverID, "DriverLocationUpdated", "driver.location."+driverID),
		DriverID:  driverID,
		Latitude:  lat,
		Longitude: lng,
	}
}

// DeliveryCompleted は配達が完了した時に発行される。
type DeliveryCompleted struct {
	events.BaseEvent
	OrderID  string `json:"order_id"`
	DriverID string `json:"driver_id"`
}

func NewDeliveryCompleted(orderID, driverID string) DeliveryCompleted {
	return DeliveryCompleted{
		BaseEvent: events.NewBaseEvent(orderID, "DeliveryCompleted", "delivery.completed"),
		OrderID:   orderID,
		DriverID:  driverID,
	}
}
