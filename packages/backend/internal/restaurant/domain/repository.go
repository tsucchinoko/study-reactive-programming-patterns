package domain

import (
	"context"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// RestaurantRepository defines the persistence interface for Restaurant aggregates.
type RestaurantRepository interface {
	Save(ctx context.Context, restaurant Restaurant) error
	FindByID(ctx context.Context, id types.RestaurantID) (Restaurant, error)
	FindAll(ctx context.Context) ([]Restaurant, error)
	FindByCuisine(ctx context.Context, cuisine CuisineType) ([]Restaurant, error)
}
