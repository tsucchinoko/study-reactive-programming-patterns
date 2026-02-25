package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/postgres"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/simulator"
	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	orderinfra "github.com/tsucchinoko/food-delivery-tracker/internal/order/infrastructure"
	restaurantinfra "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/infrastructure"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[simulator] starting...")

	// PostgreSQLに接続
	pool, err := postgres.NewPool(ctx, postgres.DefaultConfig())
	if err != nil {
		log.Fatalf("[simulator] failed to connect to database: %v", err)
	}
	defer pool.Close()

	// イベントバスのセットアップ（Phase 1ではインメモリ）
	bus := events.NewInMemoryBus()

	// リポジトリのセットアップ
	orderRepo := orderinfra.NewPostgresOrderRepository(pool)
	restaurantRepo := restaurantinfra.NewPostgresRestaurantRepository(pool)

	// サービスのセットアップ
	orderService := orderapp.NewOrderService(orderRepo, bus.Publish)

	// レストランのシードデータ投入
	restaurants := simulator.SeedRestaurants()
	for _, r := range restaurants {
		if err := restaurantRepo.Save(ctx, r); err != nil {
			log.Fatalf("[simulator] failed to seed restaurant %s: %v", r.Name(), err)
		}
	}
	log.Printf("[simulator] seeded %d restaurants", len(restaurants))

	// DBからレストランを読み込み（一貫したIDを取得するため）
	dbRestaurants, err := restaurantRepo.FindAll(ctx)
	if err != nil {
		log.Fatalf("[simulator] failed to load restaurants: %v", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ステート遷移ゴルーチンを開始
	go simulator.RunStateTransitions(ctx, rng, orderService, 0.05) // 5% cancel rate

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
