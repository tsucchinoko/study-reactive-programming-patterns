package application

import "github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"

// OpenRestaurantCommand marks a restaurant as open.
type OpenRestaurantCommand struct {
	RestaurantID types.RestaurantID
}

// CloseRestaurantCommand marks a restaurant as closed.
type CloseRestaurantCommand struct {
	RestaurantID types.RestaurantID
}
