package application

import (
	"context"
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/order/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/fp"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// OrderService orchestrates order use cases.
// Domain logic lives in domain/ — this layer handles side effects (persistence, events).
type OrderService struct {
	repo      domain.OrderRepository
	publisher events.EventPublisher
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo domain.OrderRepository, publisher events.EventPublisher) *OrderService {
	return &OrderService{repo: repo, publisher: publisher}
}

// PlaceOrder creates a new order from the command.
func (s *OrderService) PlaceOrder(ctx context.Context, cmd PlaceOrderCommand) (domain.Order, error) {
	items := fp.Map(cmd.Items, func(i PlaceOrderItem) domain.OrderItem {
		return domain.NewOrderItem(i.MenuItemID, i.Name, i.Quantity, i.UnitPrice, i.SpecialInstructions)
	})

	orderResult := domain.BuildOrder(
		types.NewOrderID(),
		cmd.CustomerID,
		cmd.RestaurantID,
		items,
		types.Now(),
	)
	if orderResult.IsErr() {
		return domain.Order{}, orderResult.UnwrapErr()
	}
	order := orderResult.Unwrap()

	if err := s.repo.Save(ctx, order); err != nil {
		return domain.Order{}, fmt.Errorf("save order: %w", err)
	}

	event := domain.NewOrderPlaced(order)
	if err := s.publisher(ctx, event); err != nil {
		return domain.Order{}, fmt.Errorf("publish OrderPlaced: %w", err)
	}

	return order, nil
}

// ConfirmOrder transitions an order to Confirmed status.
func (s *OrderService) ConfirmOrder(ctx context.Context, cmd ConfirmOrderCommand) (domain.Order, error) {
	order, err := s.repo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("find order: %w", err)
	}

	previousStatus := order.Status()
	confirmResult := order.Confirm(types.Now())
	if confirmResult.IsErr() {
		return domain.Order{}, confirmResult.UnwrapErr()
	}
	confirmed := confirmResult.Unwrap()

	if err := s.repo.Save(ctx, confirmed); err != nil {
		return domain.Order{}, fmt.Errorf("save order: %w", err)
	}

	event := domain.NewOrderStatusChanged(confirmed.ID(), previousStatus, confirmed.Status())
	if err := s.publisher(ctx, event); err != nil {
		return domain.Order{}, fmt.Errorf("publish OrderStatusChanged: %w", err)
	}

	return confirmed, nil
}

// CancelOrder transitions an order to Cancelled status.
func (s *OrderService) CancelOrder(ctx context.Context, cmd CancelOrderCommand) (domain.Order, error) {
	order, err := s.repo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("find order: %w", err)
	}

	cancelResult := order.Cancel(cmd.Reason, types.Now())
	if cancelResult.IsErr() {
		return domain.Order{}, cancelResult.UnwrapErr()
	}
	cancelled := cancelResult.Unwrap()

	if err := s.repo.Save(ctx, cancelled); err != nil {
		return domain.Order{}, fmt.Errorf("save order: %w", err)
	}

	event := domain.NewOrderCancelled(cancelled)
	if err := s.publisher(ctx, event); err != nil {
		return domain.Order{}, fmt.Errorf("publish OrderCancelled: %w", err)
	}

	return cancelled, nil
}

// TransitionOrder advances an order to the specified status.
func (s *OrderService) TransitionOrder(ctx context.Context, cmd TransitionOrderCommand) (domain.Order, error) {
	order, err := s.repo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("find order: %w", err)
	}

	previousStatus := order.Status()
	transitionResult := order.TransitionTo(domain.OrderStatus(cmd.NewStatus), types.Now())
	if transitionResult.IsErr() {
		return domain.Order{}, transitionResult.UnwrapErr()
	}
	transitioned := transitionResult.Unwrap()

	if err := s.repo.Save(ctx, transitioned); err != nil {
		return domain.Order{}, fmt.Errorf("save order: %w", err)
	}

	event := domain.NewOrderStatusChanged(transitioned.ID(), previousStatus, transitioned.Status())
	if err := s.publisher(ctx, event); err != nil {
		return domain.Order{}, fmt.Errorf("publish OrderStatusChanged: %w", err)
	}

	return transitioned, nil
}

// GetOrder retrieves a single order by ID.
func (s *OrderService) GetOrder(ctx context.Context, q GetOrderQuery) (domain.Order, error) {
	return s.repo.FindByID(ctx, q.OrderID)
}

// ListOrders retrieves orders with pagination.
func (s *OrderService) ListOrders(ctx context.Context, q ListOrdersQuery) ([]domain.Order, error) {
	return s.repo.FindAll(ctx, q.Limit, q.Offset)
}

// ListOrdersByStatus retrieves orders filtered by status.
func (s *OrderService) ListOrdersByStatus(ctx context.Context, q ListOrdersByStatusQuery) ([]domain.Order, error) {
	return s.repo.FindByStatus(ctx, domain.OrderStatus(q.Status))
}
