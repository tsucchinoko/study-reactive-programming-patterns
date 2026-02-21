package application

import "github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"

// GetRestaurantQuery は特定のレストランを取得するクエリを表す。
type GetRestaurantQuery struct {
	RestaurantID types.RestaurantID
}

// ListRestaurantsQuery は全レストランを取得するクエリを表す。
type ListRestaurantsQuery struct{}

// GetMenuQuery はレストランのメニューを取得するクエリを表す。
type GetMenuQuery struct {
	RestaurantID types.RestaurantID
}
