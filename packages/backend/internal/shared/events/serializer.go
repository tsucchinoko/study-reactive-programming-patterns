package events

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// eventEnvelope はシリアライズ時にイベント型情報を含むラッパー。
type eventEnvelope struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

var (
	registryMu sync.RWMutex
	registry   = make(map[string]reflect.Type)
)

// RegisterEventType はイベント型名と Go の型を登録する。
// アプリケーション起動時にすべてのイベント型を登録すること。
func RegisterEventType(eventType string, sample DomainEvent) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[eventType] = reflect.TypeOf(sample)
}

// SerializeEvent は DomainEvent を JSON バイト列にシリアライズする。
func SerializeEvent(event DomainEvent) ([]byte, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("serialize event payload: %w", err)
	}

	envelope := eventEnvelope{
		EventType: event.EventType(),
		Payload:   payload,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("serialize event envelope: %w", err)
	}
	return data, nil
}

// DeserializeEvent は JSON バイト列を DomainEvent に復元する。
// 登録されていないイベント型の場合は GenericEvent として復元する。
func DeserializeEvent(data []byte) (DomainEvent, error) {
	var envelope eventEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("deserialize event envelope: %w", err)
	}

	registryMu.RLock()
	typ, ok := registry[envelope.EventType]
	registryMu.RUnlock()

	if !ok {
		// 未登録の型は GenericEvent として復元
		var generic GenericEvent
		if err := json.Unmarshal(envelope.Payload, &generic); err != nil {
			return nil, fmt.Errorf("deserialize as generic event: %w", err)
		}
		return generic, nil
	}

	// ポインタ型の場合は Elem で実体を取得
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	eventPtr := reflect.New(typ).Interface()
	if err := json.Unmarshal(envelope.Payload, eventPtr); err != nil {
		return nil, fmt.Errorf("deserialize event %s: %w", envelope.EventType, err)
	}

	event, ok := reflect.ValueOf(eventPtr).Elem().Interface().(DomainEvent)
	if !ok {
		return nil, fmt.Errorf("deserialized value does not implement DomainEvent")
	}
	return event, nil
}

// GenericEvent は未登録のイベント型を汎用的に表現する。
type GenericEvent struct {
	BaseEvent
}
