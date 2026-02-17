package graphql

import (
	orderapp "github.com/daichitsuchiya/food-delivery-tracker/internal/order/application"
	restaurantapp "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/application"
)

// Resolver is the root resolver. Dependencies are injected here.
type Resolver struct {
	OrderService      *orderapp.OrderService
	RestaurantService *restaurantapp.RestaurantService
	SubscriptionMgr   *SubscriptionManager
}
