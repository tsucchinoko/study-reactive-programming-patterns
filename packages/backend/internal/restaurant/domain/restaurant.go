package domain

import (
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/result"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// CuisineType represents the type of cuisine a restaurant serves.
type CuisineType string

const (
	CuisineJapanese CuisineType = "JAPANESE"
	CuisineItalian  CuisineType = "ITALIAN"
	CuisineChinese  CuisineType = "CHINESE"
	CuisineKorean   CuisineType = "KOREAN"
	CuisineThai     CuisineType = "THAI"
	CuisineIndian   CuisineType = "INDIAN"
	CuisineAmerican CuisineType = "AMERICAN"
	CuisineFrench   CuisineType = "FRENCH"
	CuisineMexican  CuisineType = "MEXICAN"
	CuisineOther    CuisineType = "OTHER"
)

// Location is a value object representing a geographic position.
type Location struct {
	lat     float64
	lng     float64
	address string
}

// NewLocation creates a Location value object.
func NewLocation(lat, lng float64, address string) Location {
	return Location{lat: lat, lng: lng, address: address}
}

func (l Location) Lat() float64    { return l.lat }
func (l Location) Lng() float64    { return l.lng }
func (l Location) Address() string { return l.address }

// Restaurant is the aggregate root for the Restaurant bounded context.
type Restaurant struct {
	id       types.RestaurantID
	name     string
	cuisine  CuisineType
	location Location
	menu     Menu
	isOpen   bool
}

// NewRestaurant creates a new Restaurant.
func NewRestaurant(
	id types.RestaurantID,
	name string,
	cuisine CuisineType,
	location Location,
	menu Menu,
	isOpen bool,
) Restaurant {
	return Restaurant{
		id:       id,
		name:     name,
		cuisine:  cuisine,
		location: location,
		menu:     menu,
		isOpen:   isOpen,
	}
}

func (r Restaurant) ID() types.RestaurantID { return r.id }
func (r Restaurant) Name() string           { return r.name }
func (r Restaurant) Cuisine() CuisineType   { return r.cuisine }
func (r Restaurant) Location() Location     { return r.location }
func (r Restaurant) Menu() Menu             { return r.menu }
func (r Restaurant) IsOpen() bool           { return r.isOpen }

// CanAcceptOrder checks whether the restaurant can currently accept orders.
// Pure function.
func CanAcceptOrder(restaurant Restaurant) result.Result[result.Unit] {
	if !restaurant.IsOpen() {
		return result.Err[result.Unit](fmt.Errorf("restaurant %s is currently closed", restaurant.Name()))
	}
	if len(restaurant.Menu().AvailableItems()) == 0 {
		return result.Err[result.Unit](fmt.Errorf("restaurant %s has no available menu items", restaurant.Name()))
	}
	return result.OkUnit()
}

// WithOpenStatus returns a new Restaurant with updated open status.
func (r Restaurant) WithOpenStatus(isOpen bool) Restaurant {
	return Restaurant{
		id: r.id, name: r.name, cuisine: r.cuisine,
		location: r.location, menu: r.menu, isOpen: isOpen,
	}
}
