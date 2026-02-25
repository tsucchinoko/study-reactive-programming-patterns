package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	deliveryapp "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/application"
	deliverydomain "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	deliveryinfra "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/infrastructure"
	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	orderdomain "github.com/tsucchinoko/food-delivery-tracker/internal/order/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

// DeliveryFlowSimulator は注文が READY になった時に配達フローを自動化する。
// READY → ドライバーアサイン → PICKED_UP → 位置シミュレーション → DELIVERED
type DeliveryFlowSimulator struct {
	orderService    *orderapp.OrderService
	deliveryService *deliveryapp.DeliveryService
	locationSim     *deliveryinfra.LocationSimulator
	publisher       events.EventPublisher
	rng             *rand.Rand
}

// NewDeliveryFlowSimulator は新しい DeliveryFlowSimulator を生成する。
func NewDeliveryFlowSimulator(
	orderService *orderapp.OrderService,
	deliveryService *deliveryapp.DeliveryService,
	locationSim *deliveryinfra.LocationSimulator,
	publisher events.EventPublisher,
	rng *rand.Rand,
) *DeliveryFlowSimulator {
	return &DeliveryFlowSimulator{
		orderService:    orderService,
		deliveryService: deliveryService,
		locationSim:     locationSim,
		publisher:       publisher,
		rng:             rng,
	}
}

// RunDeliveryFlow は READY ステータスの注文を監視し、配達フローを開始する。
func (s *DeliveryFlowSimulator) RunDeliveryFlow(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processReadyOrders(ctx)
		}
	}
}

func (s *DeliveryFlowSimulator) processReadyOrders(ctx context.Context) {
	// READY ステータスの注文を取得
	readyOrders, err := s.orderService.ListOrdersByStatus(ctx, orderapp.ListOrdersByStatusQuery{
		Status: string(orderdomain.StatusReady),
	})
	if err != nil {
		log.Printf("[delivery-sim] error listing ready orders: %v", err)
		return
	}

	for _, order := range readyOrders {
		// 既にアサイン済みか確認
		_, err := s.deliveryService.GetAssignmentByOrder(ctx, deliveryapp.GetAssignmentByOrderQuery{
			OrderID: order.ID(),
		})
		if err == nil {
			continue // 既にアサイン済み
		}

		// ドライバーをアサイン
		// ピックアップ位置はレストランの位置（シミュレーション用に固定値を使用）
		pickupLocation := deliverydomain.NewLocation(35.6762, 139.6503) // 渋谷付近

		assignment, err := s.deliveryService.AssignDriver(ctx, deliveryapp.AssignDriverCommand{
			OrderID:        order.ID(),
			PickupLocation: pickupLocation,
		})
		if err != nil {
			log.Printf("[delivery-sim] failed to assign driver for order %s: %v", order.ID(), err)
			continue
		}

		log.Printf("[delivery-sim] 🚗 assigned driver %s to order %s",
			assignment.DriverID(), order.ID())

		// 注文を PICKED_UP に遷移
		go s.startDeliveryAfterDelay(ctx, order, assignment, pickupLocation)
	}
}

func (s *DeliveryFlowSimulator) startDeliveryAfterDelay(
	ctx context.Context,
	order orderdomain.Order,
	assignment deliverydomain.DeliveryAssignment,
	pickupLocation deliverydomain.Location,
) {
	// ピックアップまでの遅延（2-5秒）
	delay := time.Duration(2+s.rng.Intn(4)) * time.Second
	select {
	case <-ctx.Done():
		return
	case <-time.After(delay):
	}

	// PICKED_UP に遷移
	_, err := s.orderService.TransitionOrder(ctx, orderapp.TransitionOrderCommand{
		OrderID:   order.ID(),
		NewStatus: string(orderdomain.StatusPickedUp),
	})
	if err != nil {
		log.Printf("[delivery-sim] error transitioning order %s to PICKED_UP: %v", order.ID(), err)
		return
	}
	log.Printf("[delivery-sim] 📦 order %s picked up", order.ID())

	// DELIVERING に遷移
	_, err = s.orderService.TransitionOrder(ctx, orderapp.TransitionOrderCommand{
		OrderID:   order.ID(),
		NewStatus: string(orderdomain.StatusDelivering),
	})
	if err != nil {
		log.Printf("[delivery-sim] error transitioning order %s to DELIVERING: %v", order.ID(), err)
		return
	}
	log.Printf("[delivery-sim] 🛵 order %s delivering", order.ID())

	// 配達先（ランダムな東京エリアの位置）
	deliveryLocation := deliverydomain.NewLocation(
		35.67+s.rng.Float64()*0.03,
		139.68+s.rng.Float64()*0.10,
	)

	// 位置シミュレーション開始（完了時に注文を DELIVERED に遷移）
	s.locationSim.SimulateDeliveryRoute(
		ctx,
		assignment.DriverID(),
		pickupLocation,
		deliveryLocation,
		func() {
			// 配達完了
			if err := s.deliveryService.CompleteDelivery(ctx, deliveryapp.CompleteDeliveryCommand{
				OrderID: order.ID(),
			}); err != nil {
				log.Printf("[delivery-sim] error completing delivery for order %s: %v", order.ID(), err)
				return
			}

			_, err := s.orderService.TransitionOrder(ctx, orderapp.TransitionOrderCommand{
				OrderID:   order.ID(),
				NewStatus: string(orderdomain.StatusDelivered),
			})
			if err != nil {
				log.Printf("[delivery-sim] error transitioning order %s to DELIVERED: %v", order.ID(), err)
				return
			}
			log.Printf("[delivery-sim] ✅ order %s delivered", order.ID())
		},
	)
}
