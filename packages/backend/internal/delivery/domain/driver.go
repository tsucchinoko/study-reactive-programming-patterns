package domain

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/result"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// Driver は配達境界づけられたコンテキストの集約ルート。
// 全フィールドは非公開 — イミュータブル。
type Driver struct {
	id              types.DriverID
	name            string
	phone           string
	status          DriverStatus
	currentLocation Location
}

// NewDriver は Available ステータスの新しい Driver を生成する。
func NewDriver(
	id types.DriverID,
	name string,
	phone string,
	location Location,
) Driver {
	return Driver{
		id:              id,
		name:            name,
		phone:           phone,
		status:          DriverStatusAvailable,
		currentLocation: location,
	}
}

// Reconstitute は永続化されたデータから Driver を復元する。
func Reconstitute(
	id types.DriverID,
	name string,
	phone string,
	status DriverStatus,
	location Location,
) Driver {
	return Driver{
		id:              id,
		name:            name,
		phone:           phone,
		status:          status,
		currentLocation: location,
	}
}

// --- アクセサ ---

func (d Driver) ID() types.DriverID    { return d.id }
func (d Driver) Name() string           { return d.name }
func (d Driver) Phone() string          { return d.phone }
func (d Driver) Status() DriverStatus   { return d.status }
func (d Driver) CurrentLocation() Location { return d.currentLocation }

// --- 状態遷移（新しいインスタンスを返す） ---

// Assign はドライバーを Assigned ステータスに遷移させる。
func (d Driver) Assign() result.Result[Driver] {
	return result.FlatMap(
		CanDriverTransitionTo(d.status, DriverStatusAssigned),
		func(_ result.Unit) result.Result[Driver] {
			return result.Ok(Driver{
				id: d.id, name: d.name, phone: d.phone,
				status:          DriverStatusAssigned,
				currentLocation: d.currentLocation,
			})
		},
	)
}

// StartDelivery はドライバーを Delivering ステータスに遷移させる。
func (d Driver) StartDelivery() result.Result[Driver] {
	return result.FlatMap(
		CanDriverTransitionTo(d.status, DriverStatusDelivering),
		func(_ result.Unit) result.Result[Driver] {
			return result.Ok(Driver{
				id: d.id, name: d.name, phone: d.phone,
				status:          DriverStatusDelivering,
				currentLocation: d.currentLocation,
			})
		},
	)
}

// CompleteDelivery はドライバーを Available ステータスに戻す。
func (d Driver) CompleteDelivery() result.Result[Driver] {
	return result.FlatMap(
		CanDriverTransitionTo(d.status, DriverStatusAvailable),
		func(_ result.Unit) result.Result[Driver] {
			return result.Ok(Driver{
				id: d.id, name: d.name, phone: d.phone,
				status:          DriverStatusAvailable,
				currentLocation: d.currentLocation,
			})
		},
	)
}

// WithLocation は位置を更新した新しい Driver を返す。
func (d Driver) WithLocation(location Location) Driver {
	return Driver{
		id: d.id, name: d.name, phone: d.phone,
		status:          d.status,
		currentLocation: location,
	}
}
