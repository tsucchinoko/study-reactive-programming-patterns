package application

import (
	"context"
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/order/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/fp"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// OrderService は注文ユースケースをオーケストレーションする。
// ドメインロジックはdomain/にあり、このレイヤーは副作用（永続化、イベント）を処理する。
type OrderService struct {
	repo      domain.OrderRepository
	publisher events.EventPublisher
}

// NewOrderService は新しいOrderServiceを作成する。
func NewOrderService(repo domain.OrderRepository, publisher events.EventPublisher) *OrderService {
	return &OrderService{repo: repo, publisher: publisher}
}

// PlaceOrder はコマンドから新しい注文を作成する。
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

// ConfirmOrder は注文をConfirmedステータスに遷移させる。
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

// CancelOrder は注文をCancelledステータスに遷移させる。
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

// TransitionOrder は注文を指定されたステータスに進める。
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

// GetOrder はIDで単一の注文を取得する。
func (s *OrderService) GetOrder(ctx context.Context, q GetOrderQuery) (domain.Order, error) {
	return s.repo.FindByID(ctx, q.OrderID)
}

// ListOrders はページネーション付きで注文を取得する。
func (s *OrderService) ListOrders(ctx context.Context, q ListOrdersQuery) ([]domain.Order, error) {
	return s.repo.FindAll(ctx, q.Limit, q.Offset)
}

// ListOrdersByStatus はステータスでフィルタリングした注文を取得する。
func (s *OrderService) ListOrdersByStatus(ctx context.Context, q ListOrdersByStatusQuery) ([]domain.Order, error) {
	return s.repo.FindByStatus(ctx, domain.OrderStatus(q.Status))
}
