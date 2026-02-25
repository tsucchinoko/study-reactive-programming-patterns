package domain

import (
	"fmt"
	"slices"

	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/result"
)

// DriverStatus はドライバーの現在の状態を表す。
type DriverStatus string

const (
	DriverStatusAvailable  DriverStatus = "AVAILABLE"
	DriverStatusAssigned   DriverStatus = "ASSIGNED"
	DriverStatusDelivering DriverStatus = "DELIVERING"
	DriverStatusOffline    DriverStatus = "OFFLINE"
)

// validDriverTransitions は許可されたドライバーステータスの遷移を定義する。
var validDriverTransitions = map[DriverStatus][]DriverStatus{
	DriverStatusAvailable:  {DriverStatusAssigned, DriverStatusOffline},
	DriverStatusAssigned:   {DriverStatusDelivering, DriverStatusAvailable},
	DriverStatusDelivering: {DriverStatusAvailable},
	DriverStatusOffline:    {DriverStatusAvailable},
}

// CanDriverTransitionTo はドライバーステータスの遷移が有効かどうかを検証する。
// 純粋関数。
func CanDriverTransitionTo(current, next DriverStatus) result.Result[result.Unit] {
	allowed, exists := validDriverTransitions[current]
	if !exists {
		return result.Err[result.Unit](fmt.Errorf("unknown driver status: %s", current))
	}
	if slices.Contains(allowed, next) {
		return result.OkUnit()
	}
	return result.Err[result.Unit](
		fmt.Errorf("invalid driver transition: %s → %s", current, next),
	)
}

func (s DriverStatus) String() string { return string(s) }
