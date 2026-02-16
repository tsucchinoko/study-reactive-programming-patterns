package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	orderapp "github.com/daichitsuchiya/food-delivery-tracker/internal/order/application"
	orderdomain "github.com/daichitsuchiya/food-delivery-tracker/internal/order/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// transitionDelay returns a random delay for a given status transition.
// Pure function.
func transitionDelay(rng *rand.Rand, status orderdomain.OrderStatus) time.Duration {
	ranges := map[orderdomain.OrderStatus][2]int{
		orderdomain.StatusCreated:    {2, 5},
		orderdomain.StatusConfirmed:  {3, 8},
		orderdomain.StatusPreparing:  {5, 15},
		orderdomain.StatusReady:      {2, 5},
		orderdomain.StatusPickedUp:   {3, 8},
		orderdomain.StatusDelivering: {5, 15},
	}
	r, ok := ranges[status]
	if !ok {
		return 5 * time.Second
	}
	secs := r[0] + rng.Intn(r[1]-r[0]+1)
	return time.Duration(secs) * time.Second
}

// nextStatus returns the next happy-path status for simulation.
// Pure function.
func nextStatus(current orderdomain.OrderStatus) (orderdomain.OrderStatus, bool) {
	transitions := map[orderdomain.OrderStatus]orderdomain.OrderStatus{
		orderdomain.StatusCreated:    orderdomain.StatusConfirmed,
		orderdomain.StatusConfirmed:  orderdomain.StatusPreparing,
		orderdomain.StatusPreparing:  orderdomain.StatusReady,
		orderdomain.StatusReady:      orderdomain.StatusPickedUp,
		orderdomain.StatusPickedUp:   orderdomain.StatusDelivering,
		orderdomain.StatusDelivering: orderdomain.StatusDelivered,
	}
	next, ok := transitions[current]
	return next, ok
}

// RunStateTransitions periodically advances active orders through their lifecycle.
func RunStateTransitions(
	ctx context.Context,
	rng *rand.Rand,
	orderService *orderapp.OrderService,
	cancelRate float64,
) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			advanceOrders(ctx, rng, orderService, cancelRate)
		}
	}
}

func advanceOrders(
	ctx context.Context,
	rng *rand.Rand,
	orderService *orderapp.OrderService,
	cancelRate float64,
) {
	activeStatuses := []orderdomain.OrderStatus{
		orderdomain.StatusCreated,
		orderdomain.StatusConfirmed,
		orderdomain.StatusPreparing,
		orderdomain.StatusReady,
		orderdomain.StatusPickedUp,
		orderdomain.StatusDelivering,
	}

	for _, status := range activeStatuses {
		orders, err := orderService.ListOrdersByStatus(ctx, orderapp.ListOrdersByStatusQuery{
			Status: string(status),
		})
		if err != nil {
			log.Printf("[simulator] error listing orders by status %s: %v", status, err)
			continue
		}

		for _, order := range orders {
			// Random cancellation chance (only before PICKED_UP)
			if cancelRate > 0 && rng.Float64() < cancelRate && canCancel(order.Status()) {
				_, err := orderService.CancelOrder(ctx, orderapp.CancelOrderCommand{
					OrderID: order.ID(),
					Reason:  randomCancelReason(rng),
				})
				if err != nil {
					log.Printf("[simulator] error cancelling order %s: %v", order.ID(), err)
				} else {
					log.Printf("[simulator] 🚫 order %s cancelled", order.ID())
				}
				continue
			}

			// Check if enough time has passed for transition
			delay := transitionDelay(rng, order.Status())
			elapsed := order.PlacedAt().Since()
			if elapsed < delay {
				continue
			}

			next, ok := nextStatus(order.Status())
			if !ok {
				continue
			}

			var transitioned types.OrderID
			if next == orderdomain.StatusConfirmed {
				result, err := orderService.ConfirmOrder(ctx, orderapp.ConfirmOrderCommand{
					OrderID: order.ID(),
				})
				if err != nil {
					log.Printf("[simulator] error confirming order %s: %v", order.ID(), err)
					continue
				}
				transitioned = result.ID()
			} else {
				result, err := orderService.TransitionOrder(ctx, orderapp.TransitionOrderCommand{
					OrderID:   order.ID(),
					NewStatus: string(next),
				})
				if err != nil {
					log.Printf("[simulator] error transitioning order %s to %s: %v", order.ID(), next, err)
					continue
				}
				transitioned = result.ID()
			}
			log.Printf("[simulator] ✅ order %s: %s → %s", transitioned, order.Status(), next)
		}
	}
}

func canCancel(status orderdomain.OrderStatus) bool {
	return status == orderdomain.StatusCreated ||
		status == orderdomain.StatusConfirmed ||
		status == orderdomain.StatusPreparing ||
		status == orderdomain.StatusReady
}

func randomCancelReason(rng *rand.Rand) string {
	reasons := []string{
		"お客様都合のキャンセル",
		"店舗が対応できなくなった",
		"配達員が見つからない",
		"住所が不正確",
		"お客様が応答しない",
	}
	return reasons[rng.Intn(len(reasons))]
}
