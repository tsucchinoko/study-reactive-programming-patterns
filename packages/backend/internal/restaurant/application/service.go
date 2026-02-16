package application

import (
	"context"
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
)

// RestaurantService orchestrates restaurant use cases.
type RestaurantService struct {
	repo domain.RestaurantRepository
}

// NewRestaurantService creates a new RestaurantService.
func NewRestaurantService(repo domain.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

// GetRestaurant retrieves a single restaurant by ID.
func (s *RestaurantService) GetRestaurant(ctx context.Context, q GetRestaurantQuery) (domain.Restaurant, error) {
	return s.repo.FindByID(ctx, q.RestaurantID)
}

// ListRestaurants retrieves all restaurants.
func (s *RestaurantService) ListRestaurants(ctx context.Context, _ ListRestaurantsQuery) ([]domain.Restaurant, error) {
	return s.repo.FindAll(ctx)
}

// GetMenu retrieves the menu for a restaurant.
func (s *RestaurantService) GetMenu(ctx context.Context, q GetMenuQuery) (domain.Menu, error) {
	restaurant, err := s.repo.FindByID(ctx, q.RestaurantID)
	if err != nil {
		return domain.Menu{}, fmt.Errorf("find restaurant: %w", err)
	}
	return restaurant.Menu(), nil
}
