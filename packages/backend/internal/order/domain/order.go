package domain

import (
	"fmt"

	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/option"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/result"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// Order は注文境界づけられたコンテキストの集約ルート。
// 全フィールドは非公開 — 集約はイミュータブル。
// 状態変更はWith*メソッドを通じて新しいOrderインスタンスを生成する。
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

// NewOrder はCreatedステータスの新しいOrderを作成する。
// 純粋関数 — 入力を検証し合計を計算する。
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

// Reconstitute は永続化されたデータからOrderを作成する（バリデーションをスキップ）。
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

// --- アクセサ ---

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

// --- 状態遷移（新しいインスタンスを返す） ---

// Confirm は注文をConfirmedステータスに遷移させる。
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

// TransitionTo は注文を指定されたステータスに遷移させる。
// 追加データが不要なステータス用（Preparing, Ready, PickedUp, Delivering, Delivered）。
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

// Cancel は理由付きで注文をCancelledステータスに遷移させる。
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

// --- 純粋なドメイン関数 ---

// CalculateTotal は全アイテムの小計の合計を計算する。純粋関数。
func CalculateTotal(items []OrderItem) types.Money {
	total := types.JPY(0)
	for _, item := range items {
		total = total.Add(item.Subtotal())
	}
	return total
}

// copyItems はアイテムスライスの防御的コピーを返す。
func copyItems(items []OrderItem) []OrderItem {
	copied := make([]OrderItem, len(items))
	copy(copied, items)
	return copied
}
