package application

import (
	"context"
	"fmt"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
)

// RestaurantService はレストランのユースケースを統括するサービス。
type RestaurantService struct {
	repo domain.RestaurantRepository
}

// NewRestaurantService は新しい RestaurantService を作成する。
func NewRestaurantService(repo domain.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

// GetRestaurant は ID を指定して特定のレストランを取得する。
func (s *RestaurantService) GetRestaurant(ctx context.Context, q GetRestaurantQuery) (domain.Restaurant, error) {
	return s.repo.FindByID(ctx, q.RestaurantID)
}

// ListRestaurants は全レストランを取得する。
func (s *RestaurantService) ListRestaurants(ctx context.Context, _ ListRestaurantsQuery) ([]domain.Restaurant, error) {
	return s.repo.FindAll(ctx)
}

// GetMenu はレストランのメニューを取得する。
func (s *RestaurantService) GetMenu(ctx context.Context, q GetMenuQuery) (domain.Menu, error) {
	restaurant, err := s.repo.FindByID(ctx, q.RestaurantID)
	if err != nil {
		return domain.Menu{}, fmt.Errorf("find restaurant: %w", err)
	}
	return restaurant.Menu(), nil
}
