package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryapp "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/application"
	deliverydomain "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	deliveryinfra "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/infrastructure"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/postgres"
	infraredis "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/redis"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/simulator"
	infrasqs "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/sqs"
	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	orderdomain "github.com/tsucchinoko/food-delivery-tracker/internal/order/domain"
	orderinfra "github.com/tsucchinoko/food-delivery-tracker/internal/order/infrastructure"
	restaurantinfra "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/infrastructure"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[simulator] starting...")

	// イベント型をレジストリに登録
	events.RegisterEventType("OrderPlaced", orderdomain.OrderPlaced{})
	events.RegisterEventType("OrderConfirmed", orderdomain.OrderConfirmed{})
	events.RegisterEventType("OrderStatusChanged", orderdomain.OrderStatusChanged{})
	events.RegisterEventType("OrderCancelled", orderdomain.OrderCancelled{})
	events.RegisterEventType("DriverAssigned", deliverydomain.DriverAssigned{})
	events.RegisterEventType("DriverLocationUpdated", deliverydomain.DriverLocationUpdated{})
	events.RegisterEventType("DeliveryCompleted", deliverydomain.DeliveryCompleted{})

	// PostgreSQLに接続
	pool, err := postgres.NewPool(ctx, postgres.DefaultConfig())
	if err != nil {
		log.Fatalf("[simulator] failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Redis クライアントのセットアップ
	redisClient := infraredis.NewClient(infraredis.DefaultConfig())
	defer redisClient.Close()

	// SQS クライアントのセットアップ
	sqsClient := infrasqs.NewClient(infrasqs.DefaultEndpoint(), infrasqs.DefaultRegion())
	orderQueueURL, err := infrasqs.ResolveQueueURL(ctx, sqsClient, "order-events-queue")
	if err != nil {
		log.Fatalf("[simulator] failed to resolve SQS queue URL: %v", err)
	}

	// イベントパブリッシャーの合成: Redis Pub/Sub + SQS
	publisher := events.ComposePublishers(
		infraredis.NewPublisher(redisClient),
		infrasqs.NewPublisher(sqsClient, orderQueueURL),
	)

	// 位置更新用パブリッシャー（Redis のみ — SQS に位置更新を送る必要はない）
	locationPublisher := infraredis.NewPublisher(redisClient)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// リポジトリのセットアップ
	orderRepo := orderinfra.NewPostgresOrderRepository(pool)
	restaurantRepo := restaurantinfra.NewPostgresRestaurantRepository(pool)
	driverRepo := deliveryinfra.NewPostgresDriverRepository(pool)
	assignmentRepo := deliveryinfra.NewPostgresAssignmentRepository(pool)

	// サービスのセットアップ
	orderService := orderapp.NewOrderService(orderRepo, publisher)
	deliveryService := deliveryapp.NewDeliveryService(driverRepo, assignmentRepo, publisher)
	locationSim := deliveryinfra.NewLocationSimulator(locationPublisher, rng)

	// レストランのシードデータ投入
	restaurants := simulator.SeedRestaurants()
	for _, r := range restaurants {
		if err := restaurantRepo.Save(ctx, r); err != nil {
			log.Fatalf("[simulator] failed to seed restaurant %s: %v", r.Name(), err)
		}
	}
	log.Printf("[simulator] seeded %d restaurants", len(restaurants))

	// ドライバーのシードデータ投入
	drivers := simulator.SeedDrivers()
	for _, d := range drivers {
		if err := driverRepo.Save(ctx, d); err != nil {
			log.Fatalf("[simulator] failed to seed driver %s: %v", d.Name(), err)
		}
	}
	log.Printf("[simulator] seeded %d drivers", len(drivers))

	// DBからレストランを読み込み
	dbRestaurants, err := restaurantRepo.FindAll(ctx)
	if err != nil {
		log.Fatalf("[simulator] failed to load restaurants: %v", err)
	}

	// 注文ステート遷移ゴルーチンを開始（READY 手前まで自動進行、READY 以降は配達フローが処理）
	go simulator.RunStateTransitions(ctx, rng, orderService, 0.05) // 5% cancel rate

	// 配達フローシミュレーターを開始（READY → アサイン → 位置配信 → DELIVERED）
	deliveryFlow := simulator.NewDeliveryFlowSimulator(
		orderService, deliveryService, locationSim, publisher, rng,
	)
	go deliveryFlow.RunDeliveryFlow(ctx)

	// 定期的にランダムな注文を生成
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("[simulator] running — generating orders every 10s, press Ctrl+C to stop")

	for {
		select {
		case <-ctx.Done():
			log.Println("[simulator] shutting down...")
			return
		case <-ticker.C:
			cmd := simulator.GenerateRandomOrder(rng, dbRestaurants)
			order, err := orderService.PlaceOrder(ctx, cmd)
			if err != nil {
				log.Printf("[simulator] error placing order: %v", err)
				continue
			}
			log.Printf("[simulator] 🍕 new order %s at restaurant %s — total: %s",
				order.ID(), order.RestaurantID(), order.Total())
		}
	}
}
