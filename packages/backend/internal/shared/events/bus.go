package events

import "context"

// EventPublisher はドメインイベントを発行する関数型。
type EventPublisher func(ctx context.Context, event DomainEvent) error

// EventSubscriber は指定トピックのドメインイベントを受信するインターフェース。
type EventSubscriber interface {
	Subscribe(ctx context.Context, topic string, handler func(DomainEvent) error) error
}

// ComposePublishers は複数のパブリッシャーを合成する。
// すべてのパブリッシャーが各イベントを受け取り、いずれかが失敗した時点でエラーを返す。
func ComposePublishers(publishers ...EventPublisher) EventPublisher {
	return func(ctx context.Context, event DomainEvent) error {
		for _, pub := range publishers {
			if err := pub(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}
}

// InMemoryBus はプロセス内で動作するシンプルなイベントバス。
// Phase 3 で Redis Pub/Sub に置き換え予定。
type InMemoryBus struct {
	handlers map[string][]func(DomainEvent) error
}

// NewInMemoryBus は新しい InMemoryBus を生成する。
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]func(DomainEvent) error),
	}
}

// Publish はイベントのトピックに登録された全ハンドラへイベントを配信する。
func (b *InMemoryBus) Publish(_ context.Context, event DomainEvent) error {
	for _, handler := range b.handlers[event.Topic()] {
		if err := handler(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe は指定トピックにハンドラを登録する。
func (b *InMemoryBus) Subscribe(_ context.Context, topic string, handler func(DomainEvent) error) error {
	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}
