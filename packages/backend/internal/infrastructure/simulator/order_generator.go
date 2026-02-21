package simulator

import (
	"math/rand"

	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	restaurantdomain "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// GenerateRandomOrder は利用可能なレストランからランダムなPlaceOrderCommandを作成する。
// 純粋関数 — 入力と乱数ソースにのみ依存する。
func GenerateRandomOrder(rng *rand.Rand, restaurants []restaurantdomain.Restaurant) orderapp.PlaceOrderCommand {
	// ランダムなレストランを選択
	restaurant := restaurants[rng.Intn(len(restaurants))]
	availableItems := restaurant.Menu().AvailableItems()

	// 1〜3個のランダムなアイテムを選択
	numItems := rng.Intn(3) + 1
	if numItems > len(availableItems) {
		numItems = len(availableItems)
	}

	// シャッフルして先頭N個を取得
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
