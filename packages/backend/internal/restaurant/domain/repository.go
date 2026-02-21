package domain

import (
	"context"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// RestaurantRepository は Restaurant 集約の永続化インターフェースを定義する。
type RestaurantRepository interface {
	Save(ctx context.Context, restaurant Restaurant) error
	FindByID(ctx context.Context, id types.RestaurantID) (Restaurant, error)
	FindAll(ctx context.Context) ([]Restaurant, error)
	FindByCuisine(ctx context.Context, cuisine CuisineType) ([]Restaurant, error)
}
