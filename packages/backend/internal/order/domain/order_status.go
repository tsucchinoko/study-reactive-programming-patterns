package domain

import (
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
)

// OrderStatus represents the current state of an order in its lifecycle.
type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"
	StatusConfirmed  OrderStatus = "CONFIRMED"
	StatusPreparing  OrderStatus = "PREPARING"
	StatusReady      OrderStatus = "READY"
	StatusPickedUp   OrderStatus = "PICKED_UP"
	StatusDelivering OrderStatus = "DELIVERING"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

// validTransitions defines the allowed state machine transitions.
// This is a pure data declaration — no side effects.
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusCreated:    {StatusConfirmed, StatusCancelled},
	StatusConfirmed:  {StatusPreparing, StatusCancelled},
	StatusPreparing:  {StatusReady, StatusCancelled},
	StatusReady:      {StatusPickedUp, StatusCancelled},
	StatusPickedUp:   {StatusDelivering},
	StatusDelivering: {StatusDelivered},
	StatusDelivered:  {},
	StatusCancelled:  {},
}

// CanTransitionTo checks whether a status transition is valid.
// Pure function — no side effects.
func CanTransitionTo(current, next OrderStatus) result.Result[result.Unit] {
	allowed, exists := validTransitions[current]
	if !exists {
		return result.Err[result.Unit](fmt.Errorf("unknown order status: %s", current))
	}
	for _, s := range allowed {
		if s == next {
			return result.OkUnit()
		}
	}
	return result.Err[result.Unit](
		fmt.Errorf("invalid transition: %s → %s", current, next),
	)
}

// NextStatuses returns the list of valid next statuses from the current status.
// Pure function.
func NextStatuses(current OrderStatus) []OrderStatus {
	return validTransitions[current]
}

// IsTerminal returns true if no further transitions are possible.
// Pure function.
func IsTerminal(status OrderStatus) bool {
	return len(validTransitions[status]) == 0
}

// AllStatuses returns all possible order statuses.
func AllStatuses() []OrderStatus {
	return []OrderStatus{
		StatusCreated, StatusConfirmed, StatusPreparing, StatusReady,
		StatusPickedUp, StatusDelivering, StatusDelivered, StatusCancelled,
	}
}

func (s OrderStatus) String() string { return string(s) }
