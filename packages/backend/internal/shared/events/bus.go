package events

import "context"

// EventPublisher publishes a single domain event.
type EventPublisher func(ctx context.Context, event DomainEvent) error

// EventSubscriber receives domain events on a given topic.
type EventSubscriber interface {
	Subscribe(ctx context.Context, topic string, handler func(DomainEvent) error) error
}

// ComposePublishers chains multiple publishers. All publishers receive every event.
// If any publisher fails, the error is returned immediately.
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

// InMemoryBus is a simple in-process event bus for Phase 1.
// It will be replaced by Redis Pub/Sub in Phase 3.
type InMemoryBus struct {
	handlers map[string][]func(DomainEvent) error
}

// NewInMemoryBus creates a new InMemoryBus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]func(DomainEvent) error),
	}
}

// Publish sends an event to all registered handlers for its topic.
func (b *InMemoryBus) Publish(_ context.Context, event DomainEvent) error {
	for _, handler := range b.handlers[event.Topic()] {
		if err := handler(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe registers a handler for a topic.
func (b *InMemoryBus) Subscribe(_ context.Context, topic string, handler func(DomainEvent) error) error {
	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}
