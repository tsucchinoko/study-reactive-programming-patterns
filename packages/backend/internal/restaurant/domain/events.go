package domain

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
)

type PreparationStarted struct {
	events.BaseEvent
	RestaurantID string `json:"restaurant_id"`
	OrderID      string `json:"order_id"`
}

func NewPreparationStarted(restaurantID, orderID string) PreparationStarted {
	return PreparationStarted{
		BaseEvent:    events.NewBaseEvent(orderID, "PreparationStarted", "restaurant.preparation_started"),
		RestaurantID: restaurantID,
		OrderID:      orderID,
	}
}

type PreparationCompleted struct {
	events.BaseEvent
	RestaurantID string `json:"restaurant_id"`
	OrderID      string `json:"order_id"`
}

func NewPreparationCompleted(restaurantID, orderID string) PreparationCompleted {
	return PreparationCompleted{
		BaseEvent:    events.NewBaseEvent(orderID, "PreparationCompleted", "restaurant.preparation_completed"),
		RestaurantID: restaurantID,
		OrderID:      orderID,
	}
}
