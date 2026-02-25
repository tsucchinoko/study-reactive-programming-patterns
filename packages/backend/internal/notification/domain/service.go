package domain

import "fmt"

// FormatOrderNotification はイベント型と注文IDから通知メッセージを生成する。
// 純粋関数 — 副作用なし。
func FormatOrderNotification(eventType, orderID string) (NotificationType, string) {
	switch eventType {
	case "OrderPlaced":
		return NotificationTypeOrderPlaced,
			fmt.Sprintf("新しい注文 %s が作成されました", orderID)
	case "OrderConfirmed":
		return NotificationTypeOrderConfirmed,
			fmt.Sprintf("注文 %s が確認されました", orderID)
	case "OrderStatusChanged":
		return NotificationTypeOrderReady,
			fmt.Sprintf("注文 %s のステータスが変更されました", orderID)
	case "OrderCancelled":
		return NotificationTypeOrderCancelled,
			fmt.Sprintf("注文 %s がキャンセルされました", orderID)
	case "DriverAssigned":
		return NotificationTypeDriverAssigned,
			fmt.Sprintf("注文 %s にドライバーがアサインされました", orderID)
	default:
		return NotificationType(eventType),
			fmt.Sprintf("注文 %s: %s", orderID, eventType)
	}
}
