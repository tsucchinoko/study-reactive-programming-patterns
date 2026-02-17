package domain

import (
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/option"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// Order is the aggregate root for the Order bounded context.
// All fields are unexported — the aggregate is immutable.
// State changes produce new Order instances via With* methods.
type Order struct {
	id           types.OrderID
	customerID   types.CustomerID
	restaurantID types.RestaurantID
	items        []OrderItem
	status       OrderStatus
	total        types.Money
	placedAt     types.Timestamp
	confirmedAt  option.Option[types.Timestamp]
	deliveredAt  option.Option[types.Timestamp]
	cancelledAt  option.Option[types.Timestamp]
	cancelReason option.Option[string]
}

// NewOrder creates a new Order in Created status.
// Pure function — validates inputs and computes total.
func NewOrder(
	id types.OrderID,
	customerID types.CustomerID,
	restaurantID types.RestaurantID,
	items []OrderItem,
	placedAt types.Timestamp,
) result.Result[Order] {
	if len(items) == 0 {
		return result.Err[Order](fmt.Errorf("order must have at least one item"))
	}

	total := CalculateTotal(items)
	if !total.IsPositive() {
		return result.Err[Order](fmt.Errorf("order total must be positive"))
	}

	return result.Ok(Order{
		id:           id,
		customerID:   customerID,
		restaurantID: restaurantID,
		items:        copyItems(items),
		status:       StatusCreated,
		total:        total,
		placedAt:     placedAt,
		confirmedAt:  option.None[types.Timestamp](),
		deliveredAt:  option.None[types.Timestamp](),
		cancelledAt:  option.None[types.Timestamp](),
		cancelReason: option.None[string](),
	})
}

// Reconstitute creates an Order from persisted data (bypasses validation).
func Reconstitute(
	id types.OrderID,
	customerID types.CustomerID,
	restaurantID types.RestaurantID,
	items []OrderItem,
	status OrderStatus,
	total types.Money,
	placedAt types.Timestamp,
	confirmedAt option.Option[types.Timestamp],
	deliveredAt option.Option[types.Timestamp],
	cancelledAt option.Option[types.Timestamp],
	cancelReason option.Option[string],
) Order {
	return Order{
		id:           id,
		customerID:   customerID,
		restaurantID: restaurantID,
		items:        copyItems(items),
		status:       status,
		total:        total,
		placedAt:     placedAt,
		confirmedAt:  confirmedAt,
		deliveredAt:  deliveredAt,
		cancelledAt:  cancelledAt,
		cancelReason: cancelReason,
	}
}

// --- Accessors ---

func (o Order) ID() types.OrderID                { return o.id }
func (o Order) CustomerID() types.CustomerID     { return o.customerID }
func (o Order) RestaurantID() types.RestaurantID { return o.restaurantID }
func (o Order) Status() OrderStatus              { return o.status }
func (o Order) Total() types.Money               { return o.total }
func (o Order) PlacedAt() types.Timestamp        { return o.placedAt }

func (o Order) Items() []OrderItem                          { return copyItems(o.items) }
func (o Order) ConfirmedAt() option.Option[types.Timestamp] { return o.confirmedAt }
func (o Order) DeliveredAt() option.Option[types.Timestamp] { return o.deliveredAt }
func (o Order) CancelledAt() option.Option[types.Timestamp] { return o.cancelledAt }
func (o Order) CancelReason() option.Option[string]         { return o.cancelReason }

// --- State transitions (return new instances) ---

// Confirm transitions the order to Confirmed status.
func (o Order) Confirm(at types.Timestamp) result.Result[Order] {
	return result.FlatMap(
		CanTransitionTo(o.status, StatusConfirmed),
		func(_ result.Unit) result.Result[Order] {
			return result.Ok(Order{
				id: o.id, customerID: o.customerID, restaurantID: o.restaurantID,
				items: o.items, status: StatusConfirmed, total: o.total,
				placedAt: o.placedAt, confirmedAt: option.Some(at),
				deliveredAt: o.deliveredAt, cancelledAt: o.cancelledAt,
				cancelReason: o.cancelReason,
			})
		},
	)
}

// TransitionTo moves the order to the given status.
// For statuses that don't need extra data (Preparing, Ready, PickedUp, Delivering, Delivered).
func (o Order) TransitionTo(status OrderStatus, at types.Timestamp) result.Result[Order] {
	return result.FlatMap(
		CanTransitionTo(o.status, status),
		func(_ result.Unit) result.Result[Order] {
			next := Order{
				id: o.id, customerID: o.customerID, restaurantID: o.restaurantID,
				items: o.items, status: status, total: o.total,
				placedAt: o.placedAt, confirmedAt: o.confirmedAt,
				deliveredAt: o.deliveredAt, cancelledAt: o.cancelledAt,
				cancelReason: o.cancelReason,
			}
			if status == StatusDelivered {
				next.deliveredAt = option.Some(at)
			}
			return result.Ok(next)
		},
	)
}

// Cancel transitions the order to Cancelled status with a reason.
func (o Order) Cancel(reason string, at types.Timestamp) result.Result[Order] {
	return result.FlatMap(
		CanTransitionTo(o.status, StatusCancelled),
		func(_ result.Unit) result.Result[Order] {
			return result.Ok(Order{
				id: o.id, customerID: o.customerID, restaurantID: o.restaurantID,
				items: o.items, status: StatusCancelled, total: o.total,
				placedAt: o.placedAt, confirmedAt: o.confirmedAt,
				deliveredAt: o.deliveredAt, cancelledAt: option.Some(at),
				cancelReason: option.Some(reason),
			})
		},
	)
}

// --- Pure domain functions ---

// CalculateTotal computes the sum of all item subtotals. Pure function.
func CalculateTotal(items []OrderItem) types.Money {
	total := types.JPY(0)
	for _, item := range items {
		total = total.Add(item.Subtotal())
	}
	return total
}

// copyItems returns a defensive copy of the items slice.
func copyItems(items []OrderItem) []OrderItem {
	copied := make([]OrderItem, len(items))
	copy(copied, items)
	return copied
}
