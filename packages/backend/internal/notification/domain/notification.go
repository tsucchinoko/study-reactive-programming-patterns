package domain

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// NotificationType は通知の種別を表す。
type NotificationType string

const (
	NotificationTypeOrderPlaced    NotificationType = "ORDER_PLACED"
	NotificationTypeOrderConfirmed NotificationType = "ORDER_CONFIRMED"
	NotificationTypeOrderReady     NotificationType = "ORDER_READY"
	NotificationTypeOrderDelivered NotificationType = "ORDER_DELIVERED"
	NotificationTypeOrderCancelled NotificationType = "ORDER_CANCELLED"
	NotificationTypeDriverAssigned NotificationType = "DRIVER_ASSIGNED"
)

// Notification はモック通知を表す。
type Notification struct {
	id        types.NotificationID
	notifType NotificationType
	orderID   string
	message   string
	sentAt    types.Timestamp
}

// NewNotification は新しい Notification を生成する。
func NewNotification(
	notifType NotificationType,
	orderID string,
	message string,
) Notification {
	return Notification{
		id:        types.NewNotificationID(),
		notifType: notifType,
		orderID:   orderID,
		message:   message,
		sentAt:    types.Now(),
	}
}

func (n Notification) ID() types.NotificationID  { return n.id }
func (n Notification) Type() NotificationType     { return n.notifType }
func (n Notification) OrderID() string            { return n.orderID }
func (n Notification) Message() string            { return n.message }
func (n Notification) SentAt() types.Timestamp    { return n.sentAt }
