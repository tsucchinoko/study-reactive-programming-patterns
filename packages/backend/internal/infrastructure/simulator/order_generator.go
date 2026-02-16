package simulator

import (
	"math/rand"

	orderapp "github.com/daichitsuchiya/food-delivery-tracker/internal/order/application"
	restaurantdomain "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// GenerateRandomOrder creates a random PlaceOrderCommand from available restaurants.
// Pure function — only depends on its inputs and the random source.
func GenerateRandomOrder(rng *rand.Rand, restaurants []restaurantdomain.Restaurant) orderapp.PlaceOrderCommand {
	// Pick a random restaurant
	restaurant := restaurants[rng.Intn(len(restaurants))]
	availableItems := restaurant.Menu().AvailableItems()

	// Pick 1-3 random items
	numItems := rng.Intn(3) + 1
	if numItems > len(availableItems) {
		numItems = len(availableItems)
	}

	// Shuffle and take first N
	shuffled := make([]restaurantdomain.MenuItem, len(availableItems))
	copy(shuffled, availableItems)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	items := make([]orderapp.PlaceOrderItem, numItems)
	for i := 0; i < numItems; i++ {
		mi := shuffled[i]
		quantity := rng.Intn(3) + 1
		items[i] = orderapp.PlaceOrderItem{
			MenuItemID:          mi.ID(),
			Name:                mi.Name(),
			Quantity:            quantity,
			UnitPrice:           mi.Price(),
			SpecialInstructions: randomInstructions(rng),
		}
	}

	return orderapp.PlaceOrderCommand{
		CustomerID:   types.NewCustomerID(),
		RestaurantID: restaurant.ID(),
		Items:        items,
	}
}

func randomInstructions(rng *rand.Rand) string {
	if rng.Float64() < 0.7 {
		return ""
	}
	instructions := []string{
		"辛さ控えめで",
		"ネギ多めで",
		"大盛りで",
		"ご飯少なめで",
		"別添えで",
		"アレルギー: えび",
	}
	return instructions[rng.Intn(len(instructions))]
}
