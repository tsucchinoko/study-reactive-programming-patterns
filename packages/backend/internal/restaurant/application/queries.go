package application

import "github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"

// GetRestaurantQuery retrieves a single restaurant.
type GetRestaurantQuery struct {
	RestaurantID types.RestaurantID
}

// ListRestaurantsQuery retrieves all restaurants.
type ListRestaurantsQuery struct{}

// GetMenuQuery retrieves the menu for a restaurant.
type GetMenuQuery struct {
	RestaurantID types.RestaurantID
}
