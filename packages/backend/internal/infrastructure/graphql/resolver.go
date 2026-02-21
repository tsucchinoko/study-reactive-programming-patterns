package graphql

import (
	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	restaurantapp "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/application"
)

// Resolver はルートリゾルバ。依存関係をここに注入する。
type Resolver struct {
	OrderService      *orderapp.OrderService
	RestaurantService *restaurantapp.RestaurantService
	SubscriptionMgr   *SubscriptionManager
}
