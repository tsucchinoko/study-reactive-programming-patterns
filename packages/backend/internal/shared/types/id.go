package types

import "github.com/google/uuid"

// 型付き ID ラッパーは、異なる ID 型の混在をコンパイル時に防ぐ。

type OrderID struct{ value uuid.UUID }
type CustomerID struct{ value uuid.UUID }
type RestaurantID struct{ value uuid.UUID }
type MenuID struct{ value uuid.UUID }
type MenuItemID struct{ value uuid.UUID }
type DriverID struct{ value uuid.UUID }
type AssignmentID struct{ value uuid.UUID }
type NotificationID struct{ value uuid.UUID }

// NewOrderID は新しいランダムな OrderID を生成する。
func NewOrderID() OrderID { return OrderID{value: uuid.New()} }

// OrderIDFrom は既存の UUID から OrderID を生成する。
func OrderIDFrom(id uuid.UUID) OrderID { return OrderID{value: id} }

// ParseOrderID は文字列を OrderID にパースする。
func ParseOrderID(s string) (OrderID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return OrderID{}, err
	}
	return OrderID{value: id}, nil
}

func (id OrderID) String() string  { return id.value.String() }
func (id OrderID) UUID() uuid.UUID { return id.value }

// NewCustomerID は新しいランダムな CustomerID を生成する。
func NewCustomerID() CustomerID { return CustomerID{value: uuid.New()} }

// CustomerIDFrom は既存の UUID から CustomerID を生成する。
func CustomerIDFrom(id uuid.UUID) CustomerID { return CustomerID{value: id} }

// ParseCustomerID は文字列を CustomerID にパースする。
func ParseCustomerID(s string) (CustomerID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CustomerID{}, err
	}
	return CustomerID{value: id}, nil
}

func (id CustomerID) String() string  { return id.value.String() }
func (id CustomerID) UUID() uuid.UUID { return id.value }

// NewRestaurantID は新しいランダムな RestaurantID を生成する。
func NewRestaurantID() RestaurantID { return RestaurantID{value: uuid.New()} }

// RestaurantIDFrom は既存の UUID から RestaurantID を生成する。
func RestaurantIDFrom(id uuid.UUID) RestaurantID { return RestaurantID{value: id} }

// ParseRestaurantID は文字列を RestaurantID にパースする。
func ParseRestaurantID(s string) (RestaurantID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return RestaurantID{}, err
	}
	return RestaurantID{value: id}, nil
}

func (id RestaurantID) String() string  { return id.value.String() }
func (id RestaurantID) UUID() uuid.UUID { return id.value }

// NewMenuID は新しいランダムな MenuID を生成する。
func NewMenuID() MenuID { return MenuID{value: uuid.New()} }

func (id MenuID) String() string  { return id.value.String() }
func (id MenuID) UUID() uuid.UUID { return id.value }

// NewMenuItemID は新しいランダムな MenuItemID を生成する。
func NewMenuItemID() MenuItemID { return MenuItemID{value: uuid.New()} }

// MenuItemIDFrom は既存の UUID から MenuItemID を生成する。
func MenuItemIDFrom(id uuid.UUID) MenuItemID { return MenuItemID{value: id} }

func (id MenuItemID) String() string  { return id.value.String() }
func (id MenuItemID) UUID() uuid.UUID { return id.value }

// NewDriverID は新しいランダムな DriverID を生成する。
func NewDriverID() DriverID { return DriverID{value: uuid.New()} }

// DriverIDFrom は既存の UUID から DriverID を生成する。
func DriverIDFrom(id uuid.UUID) DriverID { return DriverID{value: id} }

// ParseDriverID は文字列を DriverID にパースする。
func ParseDriverID(s string) (DriverID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return DriverID{}, err
	}
	return DriverID{value: id}, nil
}

func (id DriverID) String() string  { return id.value.String() }
func (id DriverID) UUID() uuid.UUID { return id.value }

// NewAssignmentID は新しいランダムな AssignmentID を生成する。
func NewAssignmentID() AssignmentID { return AssignmentID{value: uuid.New()} }

// AssignmentIDFrom は既存の UUID から AssignmentID を生成する。
func AssignmentIDFrom(id uuid.UUID) AssignmentID { return AssignmentID{value: id} }

// ParseAssignmentID は文字列を AssignmentID にパースする。
func ParseAssignmentID(s string) (AssignmentID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return AssignmentID{}, err
	}
	return AssignmentID{value: id}, nil
}

func (id AssignmentID) String() string  { return id.value.String() }
func (id AssignmentID) UUID() uuid.UUID { return id.value }

// NewNotificationID は新しいランダムな NotificationID を生成する。
func NewNotificationID() NotificationID { return NotificationID{value: uuid.New()} }

func (id NotificationID) String() string  { return id.value.String() }
func (id NotificationID) UUID() uuid.UUID { return id.value }
