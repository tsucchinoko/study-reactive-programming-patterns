package graphql

import (
	"log"
	"sync"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql/generated"
)

// SubscriptionManager manages active GraphQL subscription channels.
// Thread-safe: multiple goroutines can subscribe/notify concurrently.
type SubscriptionManager struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]chan *generated.Order // orderID → subscriberID → channel
}

// NewSubscriptionManager creates a new SubscriptionManager.
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[string]map[string]chan *generated.Order),
	}
}

// Subscribe registers a new subscriber for the given order ID.
// Returns a read-only channel and an unsubscribe function.
func (m *SubscriptionManager) Subscribe(orderID, subscriberID string) (<-chan *generated.Order, func()) {
	ch := make(chan *generated.Order, 1)

	m.mu.Lock()
	if m.subscribers[orderID] == nil {
		m.subscribers[orderID] = make(map[string]chan *generated.Order)
	}
	m.subscribers[orderID][subscriberID] = ch
	m.mu.Unlock()

	log.Printf("[subscription] subscriber %s subscribed to order %s", subscriberID, orderID)

	unsubscribe := func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if subs, ok := m.subscribers[orderID]; ok {
			if ch, ok := subs[subscriberID]; ok {
				close(ch)
				delete(subs, subscriberID)
				log.Printf("[subscription] subscriber %s unsubscribed from order %s", subscriberID, orderID)
			}
			if len(subs) == 0 {
				delete(m.subscribers, orderID)
			}
		}
	}

	return ch, unsubscribe
}

// Notify sends the updated order to all subscribers watching the given order ID.
// Non-blocking: if a subscriber's channel is full, the update is dropped.
func (m *SubscriptionManager) Notify(orderID string, order *generated.Order) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	subs, ok := m.subscribers[orderID]
	if !ok {
		return
	}

	for subID, ch := range subs {
		select {
		case ch <- order:
			log.Printf("[subscription] notified subscriber %s for order %s (status: %s)", subID, orderID, order.Status)
		default:
			log.Printf("[subscription] dropped notification for subscriber %s (channel full)", subID)
		}
	}
}
