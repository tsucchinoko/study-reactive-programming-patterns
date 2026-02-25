package events

import (
	"context"
	"errors"
	"testing"
)

func TestInMemoryBus_PublishAndSubscribe(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	var received DomainEvent
	err := bus.Subscribe(ctx, "order.placed", func(e DomainEvent) error {
		received = e
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	if err := bus.Publish(ctx, event); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	if received == nil {
		t.Fatal("handler should have been called")
	}
	if received.AggregateID() != "order-1" {
		t.Errorf("AggregateID() = %q, want \"order-1\"", received.AggregateID())
	}
}

func TestInMemoryBus_MultipleHandlers(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	callCount := 0
	handler := func(e DomainEvent) error {
		callCount++
		return nil
	}

	bus.Subscribe(ctx, "order.placed", handler)
	bus.Subscribe(ctx, "order.placed", handler)

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	bus.Publish(ctx, event)

	if callCount != 2 {
		t.Errorf("both handlers should be called, got callCount=%d", callCount)
	}
}

func TestInMemoryBus_TopicIsolation(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	called := false
	bus.Subscribe(ctx, "order.confirmed", func(e DomainEvent) error {
		called = true
		return nil
	})

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	bus.Publish(ctx, event)

	if called {
		t.Error("handler for different topic should not be called")
	}
}

func TestInMemoryBus_HandlerError(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	expectedErr := errors.New("handler failed")
	bus.Subscribe(ctx, "order.placed", func(e DomainEvent) error {
		return expectedErr
	})

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	err := bus.Publish(ctx, event)

	if err != expectedErr {
		t.Errorf("Publish should return handler error, got %v", err)
	}
}

func TestInMemoryBus_HandlerErrorStopsExecution(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	bus.Subscribe(ctx, "order.placed", func(e DomainEvent) error {
		return errors.New("first handler fails")
	})

	secondCalled := false
	bus.Subscribe(ctx, "order.placed", func(e DomainEvent) error {
		secondCalled = true
		return nil
	})

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	bus.Publish(ctx, event)

	if secondCalled {
		t.Error("second handler should not be called after first handler fails")
	}
}

func TestInMemoryBus_NoHandlers(t *testing.T) {
	bus := NewInMemoryBus()
	ctx := context.Background()

	event := NewBaseEvent("order-1", "order.placed", "order.placed")
	err := bus.Publish(ctx, event)

	if err != nil {
		t.Errorf("Publish with no handlers should not error, got %v", err)
	}
}

func TestInMemoryBus_ImplementsEventSubscriber(t *testing.T) {
	var _ EventSubscriber = &InMemoryBus{}
}

func TestComposePublishers(t *testing.T) {
	var order []string
	pub1 := func(ctx context.Context, event DomainEvent) error {
		order = append(order, "pub1")
		return nil
	}
	pub2 := func(ctx context.Context, event DomainEvent) error {
		order = append(order, "pub2")
		return nil
	}

	composed := ComposePublishers(pub1, pub2)
	event := NewBaseEvent("agg-1", "type", "topic")
	err := composed(context.Background(), event)

	if err != nil {
		t.Fatalf("ComposePublishers returned error: %v", err)
	}
	if len(order) != 2 || order[0] != "pub1" || order[1] != "pub2" {
		t.Errorf("publishers called in wrong order: %v", order)
	}
}

func TestComposePublishers_StopsOnError(t *testing.T) {
	expectedErr := errors.New("pub1 failed")
	pub1 := func(ctx context.Context, event DomainEvent) error {
		return expectedErr
	}
	secondCalled := false
	pub2 := func(ctx context.Context, event DomainEvent) error {
		secondCalled = true
		return nil
	}

	composed := ComposePublishers(pub1, pub2)
	event := NewBaseEvent("agg-1", "type", "topic")
	err := composed(context.Background(), event)

	if err != expectedErr {
		t.Errorf("should return first error, got %v", err)
	}
	if secondCalled {
		t.Error("second publisher should not be called after first fails")
	}
}

func TestComposePublishers_Empty(t *testing.T) {
	composed := ComposePublishers()
	event := NewBaseEvent("agg-1", "type", "topic")
	err := composed(context.Background(), event)

	if err != nil {
		t.Errorf("empty ComposePublishers should not error, got %v", err)
	}
}
