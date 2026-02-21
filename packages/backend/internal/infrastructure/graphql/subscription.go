package graphql

import (
	"log"
	"sync"

	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/graphql/generated"
)

// SubscriptionManager はアクティブなGraphQLサブスクリプションチャネルを管理する。
// スレッドセーフ: 複数のゴルーチンが同時にsubscribe/notifyを呼び出せる。
type SubscriptionManager struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]chan *generated.Order // orderID → subscriberID → channel
}

// NewSubscriptionManager は新しいSubscriptionManagerを作成する。
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[string]map[string]chan *generated.Order),
	}
}

// Subscribe は指定された注文IDに新しいサブスクライバーを登録する。
// 読み取り専用チャネルとアンサブスクライブ関数を返す。
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

// Notify は指定された注文IDを監視している全サブスクライバーに更新された注文を送信する。
// ノンブロッキング: サブスクライバーのチャネルが満杯の場合、更新はドロップされる。
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
