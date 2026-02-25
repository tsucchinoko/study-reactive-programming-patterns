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

	locationMu          sync.RWMutex
	locationSubscribers map[string]map[string]chan *generated.DriverLocation // driverID → subscriberID → channel
}

// NewSubscriptionManager は新しいSubscriptionManagerを作成する。
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers:         make(map[string]map[string]chan *generated.Order),
		locationSubscribers: make(map[string]map[string]chan *generated.DriverLocation),
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

// SubscribeDriverLocation は指定されたドライバーIDの位置更新サブスクリプションを登録する。
func (m *SubscriptionManager) SubscribeDriverLocation(driverID, subscriberID string) (<-chan *generated.DriverLocation, func()) {
	ch := make(chan *generated.DriverLocation, 5) // 位置更新は高頻度なのでバッファを大きめに

	m.locationMu.Lock()
	if m.locationSubscribers[driverID] == nil {
		m.locationSubscribers[driverID] = make(map[string]chan *generated.DriverLocation)
	}
	m.locationSubscribers[driverID][subscriberID] = ch
	m.locationMu.Unlock()

	log.Printf("[subscription] subscriber %s subscribed to driver location %s", subscriberID, driverID)

	unsubscribe := func() {
		m.locationMu.Lock()
		defer m.locationMu.Unlock()
		if subs, ok := m.locationSubscribers[driverID]; ok {
			if ch, ok := subs[subscriberID]; ok {
				close(ch)
				delete(subs, subscriberID)
				log.Printf("[subscription] subscriber %s unsubscribed from driver location %s", subscriberID, driverID)
			}
			if len(subs) == 0 {
				delete(m.locationSubscribers, driverID)
			}
		}
	}

	return ch, unsubscribe
}

// NotifyDriverLocation は指定されたドライバーIDの位置更新を全サブスクライバーに送信する。
func (m *SubscriptionManager) NotifyDriverLocation(driverID string, location *generated.DriverLocation) {
	m.locationMu.RLock()
	defer m.locationMu.RUnlock()

	subs, ok := m.locationSubscribers[driverID]
	if !ok {
		return
	}

	for subID, ch := range subs {
		select {
		case ch <- location:
			log.Printf("[subscription] notified subscriber %s for driver %s location", subID, driverID)
		default:
			log.Printf("[subscription] dropped location for subscriber %s (channel full)", subID)
		}
	}
}
