package domain

import (
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// --- Order Domain Events ---

type OrderPlaced struct {
	events.BaseEvent
	CustomerID   string      `json:"customer_id"`
	RestaurantID string      `json:"restaurant_id"`
	Items        []OrderItem `json:"items"`
	Total        int64       `json:"total"`
	Currency     string      `json:"currency"`
}

func NewOrderPlaced(order Order) OrderPlaced {
	return OrderPlaced{
		BaseEvent:    events.NewBaseEvent(order.ID().String(), "OrderPlaced", "order.placed"),
		CustomerID:   order.CustomerID().String(),
		RestaurantID: order.RestaurantID().String(),
		Items:        order.Items(),
		Total:        order.Total().Amount(),
		Currency:     order.Total().Currency(),
	}
}

type OrderConfirmed struct {
	events.BaseEvent
	ConfirmedAt types.Timestamp `json:"confirmed_at"`
}

func NewOrderConfirmed(order Order) OrderConfirmed {
	return OrderConfirmed{
		BaseEvent:   events.NewBaseEvent(order.ID().String(), "OrderConfirmed", "order.confirmed"),
		ConfirmedAt: order.ConfirmedAt().Unwrap(),
	}
}

type OrderStatusChanged struct {
	events.BaseEvent
	PreviousStatus string `json:"previous_status"`
	NewStatus      string `json:"new_status"`
}

func NewOrderStatusChanged(orderID types.OrderID, previous, current OrderStatus) OrderStatusChanged {
	return OrderStatusChanged{
		BaseEvent:      events.NewBaseEvent(orderID.String(), "OrderStatusChanged", "order.status_changed"),
		PreviousStatus: previous.String(),
		NewStatus:      current.String(),
	}
}

type OrderCancelled struct {
	events.BaseEvent
	Reason      string          `json:"reason"`
	CancelledAt types.Timestamp `json:"cancelled_at"`
}

func NewOrderCancelled(order Order) OrderCancelled {
	return OrderCancelled{
		BaseEvent:   events.NewBaseEvent(order.ID().String(), "OrderCancelled", "order.cancelled"),
		Reason:      order.CancelReason().UnwrapOr(""),
		CancelledAt: order.CancelledAt().Unwrap(),
	}
}
