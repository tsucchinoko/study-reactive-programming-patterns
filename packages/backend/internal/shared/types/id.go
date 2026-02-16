package types

import "github.com/google/uuid"

// Typed ID wrappers provide compile-time safety against mixing different ID types.

type OrderID struct{ value uuid.UUID }
type CustomerID struct{ value uuid.UUID }
type RestaurantID struct{ value uuid.UUID }
type MenuID struct{ value uuid.UUID }
type MenuItemID struct{ value uuid.UUID }
type DriverID struct{ value uuid.UUID }
type AssignmentID struct{ value uuid.UUID }
type NotificationID struct{ value uuid.UUID }

// NewOrderID creates a new random OrderID.
func NewOrderID() OrderID { return OrderID{value: uuid.New()} }

// OrderIDFrom creates an OrderID from an existing UUID.
func OrderIDFrom(id uuid.UUID) OrderID { return OrderID{value: id} }

// ParseOrderID parses a string into an OrderID.
func ParseOrderID(s string) (OrderID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return OrderID{}, err
	}
	return OrderID{value: id}, nil
}

func (id OrderID) String() string  { return id.value.String() }
func (id OrderID) UUID() uuid.UUID { return id.value }

// NewCustomerID creates a new random CustomerID.
func NewCustomerID() CustomerID { return CustomerID{value: uuid.New()} }

// CustomerIDFrom creates a CustomerID from an existing UUID.
func CustomerIDFrom(id uuid.UUID) CustomerID { return CustomerID{value: id} }

// ParseCustomerID parses a string into a CustomerID.
func ParseCustomerID(s string) (CustomerID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CustomerID{}, err
	}
	return CustomerID{value: id}, nil
}

func (id CustomerID) String() string  { return id.value.String() }
func (id CustomerID) UUID() uuid.UUID { return id.value }

// NewRestaurantID creates a new random RestaurantID.
func NewRestaurantID() RestaurantID { return RestaurantID{value: uuid.New()} }

// RestaurantIDFrom creates a RestaurantID from an existing UUID.
func RestaurantIDFrom(id uuid.UUID) RestaurantID { return RestaurantID{value: id} }

// ParseRestaurantID parses a string into a RestaurantID.
func ParseRestaurantID(s string) (RestaurantID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return RestaurantID{}, err
	}
	return RestaurantID{value: id}, nil
}

func (id RestaurantID) String() string  { return id.value.String() }
func (id RestaurantID) UUID() uuid.UUID { return id.value }

// NewMenuID creates a new random MenuID.
func NewMenuID() MenuID { return MenuID{value: uuid.New()} }

func (id MenuID) String() string  { return id.value.String() }
func (id MenuID) UUID() uuid.UUID { return id.value }

// NewMenuItemID creates a new random MenuItemID.
func NewMenuItemID() MenuItemID { return MenuItemID{value: uuid.New()} }

// MenuItemIDFrom creates a MenuItemID from an existing UUID.
func MenuItemIDFrom(id uuid.UUID) MenuItemID { return MenuItemID{value: id} }

func (id MenuItemID) String() string  { return id.value.String() }
func (id MenuItemID) UUID() uuid.UUID { return id.value }

// NewDriverID creates a new random DriverID.
func NewDriverID() DriverID { return DriverID{value: uuid.New()} }

// DriverIDFrom creates a DriverID from an existing UUID.
func DriverIDFrom(id uuid.UUID) DriverID { return DriverID{value: id} }

// ParseDriverID parses a string into a DriverID.
func ParseDriverID(s string) (DriverID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return DriverID{}, err
	}
	return DriverID{value: id}, nil
}

func (id DriverID) String() string  { return id.value.String() }
func (id DriverID) UUID() uuid.UUID { return id.value }

// NewAssignmentID creates a new random AssignmentID.
func NewAssignmentID() AssignmentID { return AssignmentID{value: uuid.New()} }

func (id AssignmentID) String() string  { return id.value.String() }
func (id AssignmentID) UUID() uuid.UUID { return id.value }

// NewNotificationID creates a new random NotificationID.
func NewNotificationID() NotificationID { return NotificationID{value: uuid.New()} }

func (id NotificationID) String() string  { return id.value.String() }
func (id NotificationID) UUID() uuid.UUID { return id.value }
