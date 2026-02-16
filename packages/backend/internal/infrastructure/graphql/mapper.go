package graphql

import (
	"time"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql/generated"
	orderdomain "github.com/daichitsuchiya/food-delivery-tracker/internal/order/domain"
	restaurantdomain "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/fp"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// --- Order mapping (domain → GraphQL) ---

func toGQLOrder(o orderdomain.Order) *generated.Order {
	return &generated.Order{
		ID:           o.ID().String(),
		CustomerID:   o.CustomerID().String(),
		RestaurantID: o.RestaurantID().String(),
		Items:        fp.Map(o.Items(), toGQLOrderItem),
		Status:       generated.OrderStatus(o.Status()),
		Total:        toGQLMoney(o.Total()),
		PlacedAt:     o.PlacedAt().Time(),
		ConfirmedAt:  optionTimeToPtr(o.ConfirmedAt()),
		DeliveredAt:  optionTimeToPtr(o.DeliveredAt()),
		CancelledAt:  optionTimeToPtr(o.CancelledAt()),
		CancelReason: optionStringToPtr(o.CancelReason()),
	}
}

func toGQLOrderItem(i orderdomain.OrderItem) *generated.OrderItem {
	return &generated.OrderItem{
		MenuItemID:          i.MenuItemID().String(),
		Name:                i.Name(),
		Quantity:            i.Quantity(),
		UnitPrice:           toGQLMoney(i.UnitPrice()),
		SpecialInstructions: i.SpecialInstructions(),
		Subtotal:            toGQLMoney(i.Subtotal()),
	}
}

func toGQLOrders(orders []orderdomain.Order) []*generated.Order {
	return fp.Map(orders, func(o orderdomain.Order) *generated.Order {
		return toGQLOrder(o)
	})
}

// --- Restaurant mapping (domain → GraphQL) ---

func toGQLRestaurant(r restaurantdomain.Restaurant) *generated.Restaurant {
	return &generated.Restaurant{
		ID:      r.ID().String(),
		Name:    r.Name(),
		Cuisine: string(r.Cuisine()),
		Location: &generated.Location{
			Lat:     r.Location().Lat(),
			Lng:     r.Location().Lng(),
			Address: r.Location().Address(),
		},
		Menu: &generated.Menu{
			ID:    r.Menu().ID().String(),
			Items: fp.Map(r.Menu().Items(), toGQLMenuItem),
		},
		IsOpen: r.IsOpen(),
	}
}

func toGQLMenuItem(m restaurantdomain.MenuItem) *generated.MenuItem {
	return &generated.MenuItem{
		ID:          m.ID().String(),
		Name:        m.Name(),
		Category:    m.Category(),
		Price:       toGQLMoney(m.Price()),
		PrepTimeSec: int(m.PrepTime().Seconds()),
		Available:   m.Available(),
	}
}

func toGQLRestaurants(restaurants []restaurantdomain.Restaurant) []*generated.Restaurant {
	return fp.Map(restaurants, func(r restaurantdomain.Restaurant) *generated.Restaurant {
		return toGQLRestaurant(r)
	})
}

// --- Shared mapping ---

func toGQLMoney(m types.Money) *generated.Money {
	return &generated.Money{
		Amount:   int(m.Amount()),
		Currency: m.Currency(),
		Display:  m.String(),
	}
}

func optionTimeToPtr(o interface {
	IsSome() bool
	Unwrap() types.Timestamp
}) *time.Time {
	if !o.IsSome() {
		return nil
	}
	t := o.Unwrap().Time()
	return &t
}

func optionStringToPtr(o interface {
	IsSome() bool
	Unwrap() string
}) *string {
	if !o.IsSome() {
		return nil
	}
	s := o.Unwrap()
	return &s
}
