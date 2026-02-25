package events

import (
	"testing"
	"time"
)

func TestNewBaseEvent(t *testing.T) {
	before := time.Now()
	e := NewBaseEvent("agg-123", "order.placed", "order.placed")
	after := time.Now()

	if e.EventID() == "" {
		t.Error("EventID should not be empty")
	}
	if e.AggregateID() != "agg-123" {
		t.Errorf("AggregateID() = %q, want \"agg-123\"", e.AggregateID())
	}
	if e.EventType() != "order.placed" {
		t.Errorf("EventType() = %q, want \"order.placed\"", e.EventType())
	}
	if e.Topic() != "order.placed" {
		t.Errorf("Topic() = %q, want \"order.placed\"", e.Topic())
	}
	if e.OccurredAt().Before(before) || e.OccurredAt().After(after) {
		t.Error("OccurredAt should be around current time")
	}
}

func TestBaseEventImplementsDomainEvent(t *testing.T) {
	var _ DomainEvent = BaseEvent{}
}

func TestNewBaseEvent_UniqueIDs(t *testing.T) {
	e1 := NewBaseEvent("agg-1", "type1", "topic1")
	e2 := NewBaseEvent("agg-2", "type2", "topic2")
	if e1.EventID() == e2.EventID() {
		t.Error("different events should have different EventIDs")
	}
}
