package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent はすべてのドメインイベントが満たすべきインターフェース。
type DomainEvent interface {
	EventID() string
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
	Topic() string
}

// BaseEvent はすべてのドメインイベントに共通するフィールドを提供する。
type BaseEvent struct {
	ID          string    `json:"event_id"`
	AggregateId string    `json:"aggregate_id"`
	Type        string    `json:"event_type"`
	Timestamp   time.Time `json:"occurred_at"`
	TopicName   string    `json:"topic"`
}

// NewBaseEvent は UUID を自動生成し、現在時刻を設定した新しい BaseEvent を作成する。
func NewBaseEvent(aggregateID, eventType, topic string) BaseEvent {
	return BaseEvent{
		ID:          uuid.New().String(),
		AggregateId: aggregateID,
		Type:        eventType,
		Timestamp:   time.Now(),
		TopicName:   topic,
	}
}

func (e BaseEvent) EventID() string       { return e.ID }
func (e BaseEvent) AggregateID() string   { return e.AggregateId }
func (e BaseEvent) EventType() string     { return e.Type }
func (e BaseEvent) OccurredAt() time.Time { return e.Timestamp }
func (e BaseEvent) Topic() string         { return e.TopicName }
