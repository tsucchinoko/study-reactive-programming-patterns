package graphql

import (
	"time"

	deliverydomain "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/graphql/generated"
	orderdomain "github.com/tsucchinoko/food-delivery-tracker/internal/order/domain"
	restaurantdomain "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/fp"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// --- 注文マッピング (domain → GraphQL) ---

// ToGQLOrder はドメインのOrderをGraphQLのOrderに変換する。
func ToGQLOrder(o orderdomain.Order) *generated.Order {
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

func ToGQLOrders(orders []orderdomain.Order) []*generated.Order {
	return fp.Map(orders, func(o orderdomain.Order) *generated.Order {
		return ToGQLOrder(o)
	})
}

// --- レストランマッピング (domain → GraphQL) ---

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

// --- 配達マッピング (domain → GraphQL) ---

func toGQLDriver(d deliverydomain.Driver) *generated.Driver {
	return &generated.Driver{
		ID:    d.ID().String(),
		Name:  d.Name(),
		Phone: d.Phone(),
		Status: generated.DriverStatus(d.Status()),
		CurrentLocation: &generated.Location{
			Lat:     d.CurrentLocation().Lat(),
			Lng:     d.CurrentLocation().Lng(),
			Address: "", // ドライバーの位置にはアドレスなし
		},
	}
}

func toGQLAssignment(a deliverydomain.DeliveryAssignment, driver *deliverydomain.Driver) *generated.DeliveryAssignment {
	result := &generated.DeliveryAssignment{
		ID:         a.ID().String(),
		OrderID:    a.OrderID().String(),
		DriverID:   a.DriverID().String(),
		Status:     string(a.Status()),
		AssignedAt: a.AssignedAt().Time(),
		PickedUpAt: optionTimeToPtr(a.PickedUpAt()),
		DeliveredAt: optionTimeToPtr(a.DeliveredAt()),
	}
	if driver != nil {
		result.Driver = toGQLDriver(*driver)
	}
	return result
}

// ToGQLDriverLocation はドライバー位置イベントをGraphQL型に変換する。
func ToGQLDriverLocation(driverID string, lat, lng float64, timestamp time.Time) *generated.DriverLocation {
	return &generated.DriverLocation{
		DriverID:  driverID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: timestamp,
	}
}

// --- 共通マッピング ---

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
