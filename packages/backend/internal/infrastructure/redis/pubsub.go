package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

// NewPublisher は DomainEvent を Redis Pub/Sub チャンネルに配信する EventPublisher を返す。
// チャンネル名は event.Topic() がそのまま使われる。
func NewPublisher(client *redis.Client) events.EventPublisher {
	return func(ctx context.Context, event events.DomainEvent) error {
		data, err := events.SerializeEvent(event)
		if err != nil {
			return fmt.Errorf("redis publish: serialize: %w", err)
		}
		if err := client.Publish(ctx, event.Topic(), data).Err(); err != nil {
			return fmt.Errorf("redis publish to %s: %w", event.Topic(), err)
		}
		log.Printf("[redis] published event %s to channel %s", event.EventType(), event.Topic())
		return nil
	}
}

// Subscriber は Redis チャンネルを購読してハンドラにイベントを配信する。
type Subscriber struct {
	client *redis.Client
}

// NewSubscriber は新しい Redis Subscriber を生成する。
func NewSubscriber(client *redis.Client) *Subscriber {
	return &Subscriber{client: client}
}

// Subscribe は指定チャンネルを購読し、受信したメッセージをハンドラに渡す。
// コンテキストがキャンセルされるまでブロッキングで動作する。goroutine 内で呼び出すこと。
func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func(events.DomainEvent) error) error {
	pubsub := s.client.Subscribe(ctx, topic)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			event, err := events.DeserializeEvent([]byte(msg.Payload))
			if err != nil {
				log.Printf("[redis] failed to deserialize event on %s: %v", topic, err)
				continue
			}
			if err := handler(event); err != nil {
				log.Printf("[redis] handler error on %s: %v", topic, err)
			}
		}
	}
}

// PublishRaw は任意の JSON データを指定チャンネルに配信する。
// ドライバー位置など、DomainEvent 以外のデータに使用する。
func PublishRaw(client *redis.Client, ctx context.Context, channel string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("redis publish raw: marshal: %w", err)
	}
	return client.Publish(ctx, channel, payload).Err()
}

// SubscribeRaw は指定チャンネルの生メッセージを受信する。
// コンテキストがキャンセルされるまでブロッキングで動作する。
func SubscribeRaw(client *redis.Client, ctx context.Context, channel string, handler func([]byte) error) error {
	pubsub := client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			if err := handler([]byte(msg.Payload)); err != nil {
				log.Printf("[redis] raw handler error on %s: %v", channel, err)
			}
		}
	}
}
