package application

import (
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// GetDriverQuery はIDでドライバーを取得する。
type GetDriverQuery struct {
	DriverID types.DriverID
}

// GetAssignmentByOrderQuery は注文IDでアサインメントを取得する。
type GetAssignmentByOrderQuery struct {
	OrderID types.OrderID
}

// ListAvailableDriversQuery は利用可能なドライバーの一覧を取得する。
type ListAvailableDriversQuery struct{}
