package application

import "github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"

// OpenRestaurantCommand はレストランを営業中にするコマンドを表す。
type OpenRestaurantCommand struct {
	RestaurantID types.RestaurantID
}

// CloseRestaurantCommand はレストランを閉店状態にするコマンドを表す。
type CloseRestaurantCommand struct {
	RestaurantID types.RestaurantID
}
