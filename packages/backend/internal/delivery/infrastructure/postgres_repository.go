package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/option"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// --- PostgresDriverRepository ---

// PostgresDriverRepository は PostgreSQL を使った DriverRepository の実装。
type PostgresDriverRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresDriverRepository は新しい PostgresDriverRepository を生成する。
func NewPostgresDriverRepository(pool *pgxpool.Pool) *PostgresDriverRepository {
	return &PostgresDriverRepository{pool: pool}
}

func (r *PostgresDriverRepository) Save(ctx context.Context, driver domain.Driver) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO drivers (id, name, phone, status, latitude, longitude, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			status = EXCLUDED.status,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			updated_at = NOW()
	`,
		driver.ID().UUID(),
		driver.Name(),
		driver.Phone(),
		string(driver.Status()),
		driver.CurrentLocation().Lat(),
		driver.CurrentLocation().Lng(),
	)
	if err != nil {
		return fmt.Errorf("upsert driver: %w", err)
	}
	return nil
}

func (r *PostgresDriverRepository) FindByID(ctx context.Context, id types.DriverID) (domain.Driver, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, phone, status, latitude, longitude
		FROM drivers WHERE id = $1
	`, id.UUID())
	return scanDriver(row)
}

func (r *PostgresDriverRepository) FindByStatus(ctx context.Context, status domain.DriverStatus) ([]domain.Driver, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, phone, status, latitude, longitude
		FROM drivers WHERE status = $1
	`, string(status))
	if err != nil {
		return nil, fmt.Errorf("query drivers by status: %w", err)
	}
	defer rows.Close()
	return scanDrivers(rows)
}

func (r *PostgresDriverRepository) FindAll(ctx context.Context) ([]domain.Driver, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, phone, status, latitude, longitude
		FROM drivers ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("query all drivers: %w", err)
	}
	defer rows.Close()
	return scanDrivers(rows)
}

func scanDriver(row pgx.Row) (domain.Driver, error) {
	var id, name, phone, status string
	var lat, lng float64
	if err := row.Scan(&id, &name, &phone, &status, &lat, &lng); err != nil {
		return domain.Driver{}, fmt.Errorf("scan driver: %w", err)
	}

	driverID, err := types.ParseDriverID(id)
	if err != nil {
		return domain.Driver{}, fmt.Errorf("parse driver ID: %w", err)
	}

	return domain.Reconstitute(
		driverID,
		name,
		phone,
		domain.DriverStatus(status),
		domain.NewLocation(lat, lng),
	), nil
}

func scanDrivers(rows pgx.Rows) ([]domain.Driver, error) {
	var drivers []domain.Driver
	for rows.Next() {
		var id, name, phone, status string
		var lat, lng float64
		if err := rows.Scan(&id, &name, &phone, &status, &lat, &lng); err != nil {
			return nil, fmt.Errorf("scan driver: %w", err)
		}
		driverID, err := types.ParseDriverID(id)
		if err != nil {
			return nil, fmt.Errorf("parse driver ID: %w", err)
		}
		drivers = append(drivers, domain.Reconstitute(
			driverID, name, phone,
			domain.DriverStatus(status),
			domain.NewLocation(lat, lng),
		))
	}
	return drivers, rows.Err()
}

// --- PostgresAssignmentRepository ---

// PostgresAssignmentRepository は PostgreSQL を使った AssignmentRepository の実装。
type PostgresAssignmentRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresAssignmentRepository は新しい PostgresAssignmentRepository を生成する。
func NewPostgresAssignmentRepository(pool *pgxpool.Pool) *PostgresAssignmentRepository {
	return &PostgresAssignmentRepository{pool: pool}
}

func (r *PostgresAssignmentRepository) Save(ctx context.Context, a domain.DeliveryAssignment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO delivery_assignments (id, order_id, driver_id, status, assigned_at, picked_up_at, delivered_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			picked_up_at = EXCLUDED.picked_up_at,
			delivered_at = EXCLUDED.delivered_at,
			updated_at = NOW()
	`,
		a.ID().UUID(),
		a.OrderID().UUID(),
		a.DriverID().UUID(),
		string(a.Status()),
		a.AssignedAt().Time(),
		optionTimestampToPtr(a.PickedUpAt()),
		optionTimestampToPtr(a.DeliveredAt()),
	)
	if err != nil {
		return fmt.Errorf("upsert assignment: %w", err)
	}
	return nil
}

func (r *PostgresAssignmentRepository) FindByID(ctx context.Context, id types.AssignmentID) (domain.DeliveryAssignment, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, order_id, driver_id, status, assigned_at, picked_up_at, delivered_at
		FROM delivery_assignments WHERE id = $1
	`, id.UUID())
	return scanAssignment(row)
}

func (r *PostgresAssignmentRepository) FindByOrderID(ctx context.Context, orderID types.OrderID) (domain.DeliveryAssignment, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, order_id, driver_id, status, assigned_at, picked_up_at, delivered_at
		FROM delivery_assignments WHERE order_id = $1
		ORDER BY assigned_at DESC LIMIT 1
	`, orderID.UUID())
	return scanAssignment(row)
}

func (r *PostgresAssignmentRepository) FindByDriverID(ctx context.Context, driverID types.DriverID) ([]domain.DeliveryAssignment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, order_id, driver_id, status, assigned_at, picked_up_at, delivered_at
		FROM delivery_assignments WHERE driver_id = $1
		ORDER BY assigned_at DESC
	`, driverID.UUID())
	if err != nil {
		return nil, fmt.Errorf("query assignments by driver: %w", err)
	}
	defer rows.Close()

	var assignments []domain.DeliveryAssignment
	for rows.Next() {
		a, err := scanAssignmentFromRows(rows)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func scanAssignment(row pgx.Row) (domain.DeliveryAssignment, error) {
	var id, orderID, driverID, status string
	var assignedAt time.Time
	var pickedUpAt, deliveredAt *time.Time

	if err := row.Scan(&id, &orderID, &driverID, &status, &assignedAt, &pickedUpAt, &deliveredAt); err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("scan assignment: %w", err)
	}

	return reconstituteAssignment(id, orderID, driverID, status, assignedAt, pickedUpAt, deliveredAt)
}

func scanAssignmentFromRows(rows pgx.Rows) (domain.DeliveryAssignment, error) {
	var id, orderID, driverID, status string
	var assignedAt time.Time
	var pickedUpAt, deliveredAt *time.Time

	if err := rows.Scan(&id, &orderID, &driverID, &status, &assignedAt, &pickedUpAt, &deliveredAt); err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("scan assignment: %w", err)
	}

	return reconstituteAssignment(id, orderID, driverID, status, assignedAt, pickedUpAt, deliveredAt)
}

func reconstituteAssignment(id, orderID, driverID, status string, assignedAt time.Time, pickedUpAt, deliveredAt *time.Time) (domain.DeliveryAssignment, error) {
	aID, err := types.ParseAssignmentID(id)
	if err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("parse assignment ID: %w", err)
	}
	oID, err := types.ParseOrderID(orderID)
	if err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("parse order ID: %w", err)
	}
	dID, err := types.ParseDriverID(driverID)
	if err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("parse driver ID: %w", err)
	}

	return domain.ReconstituteAssignment(
		aID,
		oID,
		dID,
		domain.AssignmentStatus(status),
		types.TimestampFrom(assignedAt),
		optionTimestampFromPtr(pickedUpAt),
		optionTimestampFromPtr(deliveredAt),
	), nil
}

// --- ヘルパー関数 ---

func optionTimestampToPtr(opt option.Option[types.Timestamp]) *time.Time {
	if opt.IsNone() {
		return nil
	}
	t := opt.Unwrap().Time()
	return &t
}

func optionTimestampFromPtr(t *time.Time) option.Option[types.Timestamp] {
	if t == nil {
		return option.None[types.Timestamp]()
	}
	return option.Some(types.TimestampFrom(*t))
}
