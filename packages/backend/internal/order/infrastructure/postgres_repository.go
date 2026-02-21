package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/daichitsuchiya/food-delivery-tracker/internal/order/domain"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/option"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

// PostgresOrderRepository はPostgreSQLを使ったdomain.OrderRepositoryの実装。
type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresOrderRepository は新しいPostgresOrderRepositoryを作成する。
func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{pool: pool}
}

func (r *PostgresOrderRepository) Save(ctx context.Context, order domain.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, customer_id, restaurant_id, status, total_amount, currency, placed_at, confirmed_at, delivered_at, cancelled_at, cancel_reason, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			confirmed_at = EXCLUDED.confirmed_at,
			delivered_at = EXCLUDED.delivered_at,
			cancelled_at = EXCLUDED.cancelled_at,
			cancel_reason = EXCLUDED.cancel_reason,
			updated_at = NOW()
	`,
		order.ID().UUID(),
		order.CustomerID().UUID(),
		order.RestaurantID().UUID(),
		string(order.Status()),
		order.Total().Amount(),
		order.Total().Currency(),
		order.PlacedAt().Time(),
		optionToPtr(order.ConfirmedAt()),
		optionToPtr(order.DeliveredAt()),
		optionToPtr(order.CancelledAt()),
		optionStringToPtr(order.CancelReason()),
	)
	if err != nil {
		return fmt.Errorf("upsert order: %w", err)
	}

	// 既存のアイテムを削除して再挿入（upsertのシンプルなアプローチ）
	_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, order.ID().UUID())
	if err != nil {
		return fmt.Errorf("delete order items: %w", err)
	}

	for _, item := range order.Items() {
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (order_id, menu_item_id, name, quantity, unit_price, currency, special_instructions)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
			order.ID().UUID(),
			item.MenuItemID().UUID(),
			item.Name(),
			item.Quantity(),
			item.UnitPrice().Amount(),
			item.UnitPrice().Currency(),
			item.SpecialInstructions(),
		)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id types.OrderID) (domain.Order, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, customer_id, restaurant_id, status, total_amount, currency,
		       placed_at, confirmed_at, delivered_at, cancelled_at, cancel_reason
		FROM orders WHERE id = $1
	`, id.UUID())

	order, err := scanOrder(row)
	if err != nil {
		return domain.Order{}, fmt.Errorf("find order %s: %w", id, err)
	}

	items, err := r.findOrderItems(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}

	return reconstitute(order, items), nil
}

func (r *PostgresOrderRepository) FindByCustomerID(ctx context.Context, customerID types.CustomerID) ([]domain.Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, customer_id, restaurant_id, status, total_amount, currency,
		       placed_at, confirmed_at, delivered_at, cancelled_at, cancel_reason
		FROM orders WHERE customer_id = $1 ORDER BY placed_at DESC
	`, customerID.UUID())
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	return r.scanOrders(ctx, rows)
}

func (r *PostgresOrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus) ([]domain.Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, customer_id, restaurant_id, status, total_amount, currency,
		       placed_at, confirmed_at, delivered_at, cancelled_at, cancel_reason
		FROM orders WHERE status = $1 ORDER BY placed_at DESC
	`, string(status))
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	return r.scanOrders(ctx, rows)
}

func (r *PostgresOrderRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, customer_id, restaurant_id, status, total_amount, currency,
		       placed_at, confirmed_at, delivered_at, cancelled_at, cancel_reason
		FROM orders ORDER BY placed_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	return r.scanOrders(ctx, rows)
}

// --- 内部ヘルパー ---

type orderRow struct {
	id           string
	customerID   string
	restaurantID string
	status       string
	totalAmount  int64
	currency     string
	placedAt     time.Time
	confirmedAt  *time.Time
	deliveredAt  *time.Time
	cancelledAt  *time.Time
	cancelReason *string
}

func scanOrder(row pgx.Row) (orderRow, error) {
	var o orderRow
	err := row.Scan(
		&o.id, &o.customerID, &o.restaurantID, &o.status,
		&o.totalAmount, &o.currency, &o.placedAt,
		&o.confirmedAt, &o.deliveredAt, &o.cancelledAt, &o.cancelReason,
	)
	return o, err
}

func reconstitute(o orderRow, items []domain.OrderItem) domain.Order {
	orderID, _ := types.ParseOrderID(o.id)
	customerID, _ := types.ParseCustomerID(o.customerID)
	restaurantID, _ := types.ParseRestaurantID(o.restaurantID)

	return domain.Reconstitute(
		orderID,
		customerID,
		restaurantID,
		items,
		domain.OrderStatus(o.status),
		types.NewMoney(o.totalAmount, o.currency),
		types.TimestampFrom(o.placedAt),
		timeToOption(o.confirmedAt),
		timeToOption(o.deliveredAt),
		timeToOption(o.cancelledAt),
		ptrToOption(o.cancelReason),
	)
}

func (r *PostgresOrderRepository) findOrderItems(ctx context.Context, orderID types.OrderID) ([]domain.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT menu_item_id, name, quantity, unit_price, currency, special_instructions
		FROM order_items WHERE order_id = $1
	`, orderID.UUID())
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var menuItemID string
		var name string
		var quantity int
		var unitPrice int64
		var currency string
		var instructions string
		if err := rows.Scan(&menuItemID, &name, &quantity, &unitPrice, &currency, &instructions); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		mid, _ := types.ParseOrderID(menuItemID) // reuse UUID parse
		items = append(items, domain.NewOrderItem(
			types.MenuItemIDFrom(mid.UUID()),
			name, quantity,
			types.NewMoney(unitPrice, currency),
			instructions,
		))
	}
	return items, rows.Err()
}

func (r *PostgresOrderRepository) scanOrders(ctx context.Context, rows pgx.Rows) ([]domain.Order, error) {
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o orderRow
		err := rows.Scan(
			&o.id, &o.customerID, &o.restaurantID, &o.status,
			&o.totalAmount, &o.currency, &o.placedAt,
			&o.confirmedAt, &o.deliveredAt, &o.cancelledAt, &o.cancelReason,
		)
		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orderID, _ := types.ParseOrderID(o.id)
		items, err := r.findOrderItems(ctx, orderID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, reconstitute(o, items))
	}
	return orders, rows.Err()
}

func optionToPtr(o option.Option[types.Timestamp]) *time.Time {
	if o.IsNone() {
		return nil
	}
	t := o.Unwrap().Time()
	return &t
}

func optionStringToPtr(o option.Option[string]) *string {
	if o.IsNone() {
		return nil
	}
	s := o.Unwrap()
	return &s
}

func timeToOption(t *time.Time) option.Option[types.Timestamp] {
	if t == nil {
		return option.None[types.Timestamp]()
	}
	return option.Some(types.TimestampFrom(*t))
}

func ptrToOption[T any](p *T) option.Option[T] {
	if p == nil {
		return option.None[T]()
	}
	return option.Some(*p)
}
