package domain

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/option"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// AssignmentStatus は配達アサインメントの状態を表す。
type AssignmentStatus string

const (
	AssignmentStatusAssigned   AssignmentStatus = "ASSIGNED"
	AssignmentStatusPickedUp   AssignmentStatus = "PICKED_UP"
	AssignmentStatusDelivering AssignmentStatus = "DELIVERING"
	AssignmentStatusDelivered  AssignmentStatus = "DELIVERED"
	AssignmentStatusFailed     AssignmentStatus = "FAILED"
)

// DeliveryAssignment は注文とドライバーの紐づけを表すエンティティ。
type DeliveryAssignment struct {
	id          types.AssignmentID
	orderID     types.OrderID
	driverID    types.DriverID
	status      AssignmentStatus
	assignedAt  types.Timestamp
	pickedUpAt  option.Option[types.Timestamp]
	deliveredAt option.Option[types.Timestamp]
}

// NewDeliveryAssignment は新しい配達アサインメントを生成する。
func NewDeliveryAssignment(
	id types.AssignmentID,
	orderID types.OrderID,
	driverID types.DriverID,
	assignedAt types.Timestamp,
) DeliveryAssignment {
	return DeliveryAssignment{
		id:          id,
		orderID:     orderID,
		driverID:    driverID,
		status:      AssignmentStatusAssigned,
		assignedAt:  assignedAt,
		pickedUpAt:  option.None[types.Timestamp](),
		deliveredAt: option.None[types.Timestamp](),
	}
}

// ReconstituteAssignment は永続化されたデータから DeliveryAssignment を復元する。
func ReconstituteAssignment(
	id types.AssignmentID,
	orderID types.OrderID,
	driverID types.DriverID,
	status AssignmentStatus,
	assignedAt types.Timestamp,
	pickedUpAt option.Option[types.Timestamp],
	deliveredAt option.Option[types.Timestamp],
) DeliveryAssignment {
	return DeliveryAssignment{
		id:          id,
		orderID:     orderID,
		driverID:    driverID,
		status:      status,
		assignedAt:  assignedAt,
		pickedUpAt:  pickedUpAt,
		deliveredAt: deliveredAt,
	}
}

// --- アクセサ ---

func (a DeliveryAssignment) ID() types.AssignmentID                    { return a.id }
func (a DeliveryAssignment) OrderID() types.OrderID                    { return a.orderID }
func (a DeliveryAssignment) DriverID() types.DriverID                  { return a.driverID }
func (a DeliveryAssignment) Status() AssignmentStatus                  { return a.status }
func (a DeliveryAssignment) AssignedAt() types.Timestamp               { return a.assignedAt }
func (a DeliveryAssignment) PickedUpAt() option.Option[types.Timestamp] { return a.pickedUpAt }
func (a DeliveryAssignment) DeliveredAt() option.Option[types.Timestamp] { return a.deliveredAt }

// --- 状態遷移 ---

// PickUp はアサインメントを PickedUp ステータスに遷移させる。
func (a DeliveryAssignment) PickUp(at types.Timestamp) DeliveryAssignment {
	return DeliveryAssignment{
		id: a.id, orderID: a.orderID, driverID: a.driverID,
		status:      AssignmentStatusPickedUp,
		assignedAt:  a.assignedAt,
		pickedUpAt:  option.Some(at),
		deliveredAt: a.deliveredAt,
	}
}

// StartDelivering はアサインメントを Delivering ステータスに遷移させる。
func (a DeliveryAssignment) StartDelivering() DeliveryAssignment {
	return DeliveryAssignment{
		id: a.id, orderID: a.orderID, driverID: a.driverID,
		status:      AssignmentStatusDelivering,
		assignedAt:  a.assignedAt,
		pickedUpAt:  a.pickedUpAt,
		deliveredAt: a.deliveredAt,
	}
}

// Complete はアサインメントを Delivered ステータスに遷移させる。
func (a DeliveryAssignment) Complete(at types.Timestamp) DeliveryAssignment {
	return DeliveryAssignment{
		id: a.id, orderID: a.orderID, driverID: a.driverID,
		status:      AssignmentStatusDelivered,
		assignedAt:  a.assignedAt,
		pickedUpAt:  a.pickedUpAt,
		deliveredAt: option.Some(at),
	}
}

func (s AssignmentStatus) String() string { return string(s) }
