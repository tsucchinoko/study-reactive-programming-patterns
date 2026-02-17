package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// PostgresRestaurantRepository implements domain.RestaurantRepository using PostgreSQL.
type PostgresRestaurantRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRestaurantRepository creates a new PostgresRestaurantRepository.
func NewPostgresRestaurantRepository(pool *pgxpool.Pool) *PostgresRestaurantRepository {
	return &PostgresRestaurantRepository{pool: pool}
}

func (r *PostgresRestaurantRepository) Save(ctx context.Context, restaurant domain.Restaurant) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO restaurants (id, name, cuisine, lat, lng, address, is_open)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			cuisine = EXCLUDED.cuisine,
			lat = EXCLUDED.lat,
			lng = EXCLUDED.lng,
			address = EXCLUDED.address,
			is_open = EXCLUDED.is_open
	`,
		restaurant.ID().UUID(),
		restaurant.Name(),
		string(restaurant.Cuisine()),
		restaurant.Location().Lat(),
		restaurant.Location().Lng(),
		restaurant.Location().Address(),
		restaurant.IsOpen(),
	)
	if err != nil {
		return fmt.Errorf("upsert restaurant: %w", err)
	}

	// Re-insert menu items
	_, err = tx.Exec(ctx, `DELETE FROM menu_items WHERE restaurant_id = $1`, restaurant.ID().UUID())
	if err != nil {
		return fmt.Errorf("delete menu items: %w", err)
	}

	for _, item := range restaurant.Menu().Items() {
		_, err = tx.Exec(ctx, `
			INSERT INTO menu_items (id, restaurant_id, name, category, price, currency, prep_time_sec, available)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
			item.ID().UUID(),
			restaurant.ID().UUID(),
			item.Name(),
			item.Category(),
			item.Price().Amount(),
			item.Price().Currency(),
			int(item.PrepTime().Seconds()),
			item.Available(),
		)
		if err != nil {
			return fmt.Errorf("insert menu item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRestaurantRepository) FindByID(ctx context.Context, id types.RestaurantID) (domain.Restaurant, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, cuisine, lat, lng, address, is_open
		FROM restaurants WHERE id = $1
	`, id.UUID())

	rest, err := scanRestaurant(row)
	if err != nil {
		return domain.Restaurant{}, fmt.Errorf("find restaurant %s: %w", id, err)
	}

	items, err := r.findMenuItems(ctx, id)
	if err != nil {
		return domain.Restaurant{}, err
	}

	return reconstituteRestaurant(rest, items), nil
}

func (r *PostgresRestaurantRepository) FindAll(ctx context.Context) ([]domain.Restaurant, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, cuisine, lat, lng, address, is_open
		FROM restaurants ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("query restaurants: %w", err)
	}
	return r.scanRestaurants(ctx, rows)
}

func (r *PostgresRestaurantRepository) FindByCuisine(ctx context.Context, cuisine domain.CuisineType) ([]domain.Restaurant, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, cuisine, lat, lng, address, is_open
		FROM restaurants WHERE cuisine = $1 ORDER BY name
	`, string(cuisine))
	if err != nil {
		return nil, fmt.Errorf("query restaurants: %w", err)
	}
	return r.scanRestaurants(ctx, rows)
}

// --- internal helpers ---

type restaurantRow struct {
	id      string
	name    string
	cuisine string
	lat     float64
	lng     float64
	address string
	isOpen  bool
}

func scanRestaurant(row pgx.Row) (restaurantRow, error) {
	var r restaurantRow
	err := row.Scan(&r.id, &r.name, &r.cuisine, &r.lat, &r.lng, &r.address, &r.isOpen)
	return r, err
}

func reconstituteRestaurant(r restaurantRow, items []domain.MenuItem) domain.Restaurant {
	restID, _ := types.ParseRestaurantID(r.id)
	menu := domain.NewMenu(types.NewMenuID(), items)
	location := domain.NewLocation(r.lat, r.lng, r.address)
	return domain.NewRestaurant(restID, r.name, domain.CuisineType(r.cuisine), location, menu, r.isOpen)
}

func (r *PostgresRestaurantRepository) findMenuItems(ctx context.Context, restaurantID types.RestaurantID) ([]domain.MenuItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, category, price, currency, prep_time_sec, available
		FROM menu_items WHERE restaurant_id = $1 ORDER BY category, name
	`, restaurantID.UUID())
	if err != nil {
		return nil, fmt.Errorf("query menu items: %w", err)
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var id string
		var name, category, currency string
		var price int64
		var prepTimeSec int
		var available bool
		if err := rows.Scan(&id, &name, &category, &price, &currency, &prepTimeSec, &available); err != nil {
			return nil, fmt.Errorf("scan menu item: %w", err)
		}
		mid, _ := types.ParseRestaurantID(id) // reuse UUID parse
		items = append(items, domain.NewMenuItem(
			types.MenuItemIDFrom(mid.UUID()),
			name, category,
			types.NewMoney(price, currency),
			time.Duration(prepTimeSec)*time.Second,
			available,
		))
	}
	return items, rows.Err()
}

func (r *PostgresRestaurantRepository) scanRestaurants(ctx context.Context, rows pgx.Rows) ([]domain.Restaurant, error) {
	defer rows.Close()

	var restaurants []domain.Restaurant
	for rows.Next() {
		var row restaurantRow
		err := rows.Scan(&row.id, &row.name, &row.cuisine, &row.lat, &row.lng, &row.address, &row.isOpen)
		if err != nil {
			return nil, fmt.Errorf("scan restaurant: %w", err)
		}
		restID, _ := types.ParseRestaurantID(row.id)
		items, err := r.findMenuItems(ctx, restID)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, reconstituteRestaurant(row, items))
	}
	return restaurants, rows.Err()
}
