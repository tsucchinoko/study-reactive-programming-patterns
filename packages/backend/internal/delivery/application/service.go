package application

import (
	"context"
	"fmt"

	"github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// DeliveryService は配達ユースケースをオーケストレーションする。
type DeliveryService struct {
	driverRepo     domain.DriverRepository
	assignmentRepo domain.AssignmentRepository
	publisher      events.EventPublisher
}

// NewDeliveryService は新しい DeliveryService を生成する。
func NewDeliveryService(
	driverRepo domain.DriverRepository,
	assignmentRepo domain.AssignmentRepository,
	publisher events.EventPublisher,
) *DeliveryService {
	return &DeliveryService{
		driverRepo:     driverRepo,
		assignmentRepo: assignmentRepo,
		publisher:      publisher,
	}
}

// AssignDriver は利用可能なドライバーの中から最適なドライバーを選択し、注文にアサインする。
func (s *DeliveryService) AssignDriver(ctx context.Context, cmd AssignDriverCommand) (domain.DeliveryAssignment, error) {
	// 利用可能なドライバーを検索
	available, err := s.driverRepo.FindByStatus(ctx, domain.DriverStatusAvailable)
	if err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("find available drivers: %w", err)
	}

	// 最適なドライバーを選択（純粋関数）
	bestDriver := domain.SelectBestDriver(available, cmd.PickupLocation)
	if bestDriver.IsNone() {
		return domain.DeliveryAssignment{}, fmt.Errorf("no available drivers")
	}
	driver := bestDriver.Unwrap()

	// ドライバーをアサイン状態に遷移
	assignedResult := driver.Assign()
	if assignedResult.IsErr() {
		return domain.DeliveryAssignment{}, assignedResult.UnwrapErr()
	}
	assignedDriver := assignedResult.Unwrap()

	if err := s.driverRepo.Save(ctx, assignedDriver); err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("save driver: %w", err)
	}

	// アサインメントを作成
	assignment := domain.NewDeliveryAssignment(
		types.NewAssignmentID(),
		cmd.OrderID,
		assignedDriver.ID(),
		types.Now(),
	)

	if err := s.assignmentRepo.Save(ctx, assignment); err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("save assignment: %w", err)
	}

	// イベント発行
	event := domain.NewDriverAssigned(cmd.OrderID.String(), assignedDriver.ID().String())
	if err := s.publisher(ctx, event); err != nil {
		return domain.DeliveryAssignment{}, fmt.Errorf("publish DriverAssigned: %w", err)
	}

	return assignment, nil
}

// UpdateLocation はドライバーの位置を更新し、Redis Pub/Sub で配信する。
func (s *DeliveryService) UpdateLocation(ctx context.Context, cmd UpdateLocationCommand) error {
	driver, err := s.driverRepo.FindByID(ctx, cmd.DriverID)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}

	updated := driver.WithLocation(cmd.Location)
	if err := s.driverRepo.Save(ctx, updated); err != nil {
		return fmt.Errorf("save driver location: %w", err)
	}

	// 位置更新イベントは driver.location.{driverId} チャンネルに配信
	event := domain.NewDriverLocationUpdated(
		cmd.DriverID.String(),
		cmd.Location.Lat(),
		cmd.Location.Lng(),
	)
	if err := s.publisher(ctx, event); err != nil {
		return fmt.Errorf("publish DriverLocationUpdated: %w", err)
	}

	return nil
}

// CompleteDelivery は配達を完了し、ドライバーを利用可能状態に戻す。
func (s *DeliveryService) CompleteDelivery(ctx context.Context, cmd CompleteDeliveryCommand) error {
	assignment, err := s.assignmentRepo.FindByOrderID(ctx, cmd.OrderID)
	if err != nil {
		return fmt.Errorf("find assignment: %w", err)
	}

	// アサインメントを完了に遷移
	completed := assignment.Complete(types.Now())
	if err := s.assignmentRepo.Save(ctx, completed); err != nil {
		return fmt.Errorf("save assignment: %w", err)
	}

	// ドライバーを利用可能に戻す
	driver, err := s.driverRepo.FindByID(ctx, assignment.DriverID())
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}

	completedResult := driver.CompleteDelivery()
	if completedResult.IsErr() {
		return completedResult.UnwrapErr()
	}
	availableDriver := completedResult.Unwrap()

	if err := s.driverRepo.Save(ctx, availableDriver); err != nil {
		return fmt.Errorf("save driver: %w", err)
	}

	// イベント発行
	event := domain.NewDeliveryCompleted(cmd.OrderID.String(), assignment.DriverID().String())
	if err := s.publisher(ctx, event); err != nil {
		return fmt.Errorf("publish DeliveryCompleted: %w", err)
	}

	return nil
}

// GetDriver はIDでドライバーを取得する。
func (s *DeliveryService) GetDriver(ctx context.Context, q GetDriverQuery) (domain.Driver, error) {
	return s.driverRepo.FindByID(ctx, q.DriverID)
}

// GetAssignmentByOrder は注文IDでアサインメントを取得する。
func (s *DeliveryService) GetAssignmentByOrder(ctx context.Context, q GetAssignmentByOrderQuery) (domain.DeliveryAssignment, error) {
	return s.assignmentRepo.FindByOrderID(ctx, q.OrderID)
}

// ListAvailableDrivers は利用可能なドライバーの一覧を取得する。
func (s *DeliveryService) ListAvailableDrivers(ctx context.Context) ([]domain.Driver, error) {
	return s.driverRepo.FindByStatus(ctx, domain.DriverStatusAvailable)
}
