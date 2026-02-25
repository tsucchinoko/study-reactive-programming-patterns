package application

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tsucchinoko/food-delivery-tracker/internal/notification/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

// OrderEventHandler は SQS から受信した注文イベントを処理し、モック通知を送信する。
func OrderEventHandler(_ context.Context, body []byte) error {
	// エンベロープからイベント型と集約IDを抽出
	var envelope struct {
		EventType string          `json:"event_type"`
		Payload   json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		log.Printf("[notification] failed to parse event envelope: %v", err)
		return nil // パースエラーはリトライしない
	}

	// ペイロードから集約IDを抽出
	var base events.BaseEvent
	if err := json.Unmarshal(envelope.Payload, &base); err != nil {
		log.Printf("[notification] failed to parse event payload: %v", err)
		return nil
	}

	// 純粋関数で通知メッセージを生成
	notifType, message := domain.FormatOrderNotification(envelope.EventType, base.AggregateID())

	// モック通知を生成して「送信」（ログ出力）
	notification := domain.NewNotification(notifType, base.AggregateID(), message)
	log.Printf("[notification] 📧 %s | %s | %s",
		notification.Type(),
		notification.OrderID(),
		notification.Message(),
	)

	return nil
}
