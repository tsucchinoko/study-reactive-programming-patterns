package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gql "github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/postgres"
	orderapp "github.com/daichitsuchiya/food-delivery-tracker/internal/order/application"
	orderinfra "github.com/daichitsuchiya/food-delivery-tracker/internal/order/infrastructure"
	restaurantapp "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/application"
	restaurantinfra "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/infrastructure"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/types"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[api] starting...")

	// PostgreSQLに接続
	pool, err := postgres.NewPool(ctx, postgres.DefaultConfig())
	if err != nil {
		log.Fatalf("[api] failed to connect to database: %v", err)
	}
	defer pool.Close()

	// イベントバスのセットアップ（Phase 1ではインメモリ）
	bus := events.NewInMemoryBus()

	// GraphQL Subscription用のサブスクリプションマネージャーをセットアップ
	subMgr := gql.NewSubscriptionManager()

	// リポジトリのセットアップ
	orderRepo := orderinfra.NewPostgresOrderRepository(pool)
	restaurantRepo := restaurantinfra.NewPostgresRestaurantRepository(pool)

	// サービスのセットアップ
	orderService := orderapp.NewOrderService(orderRepo, bus.Publish)
	restaurantService := restaurantapp.NewRestaurantService(restaurantRepo)

	// イベントバス → サブスクリプションマネージャーの接続
	// 注文イベント発火時に注文を再取得し、サブスクライバーに通知する
	orderTopics := []string{"order.placed", "order.confirmed", "order.status_changed", "order.cancelled"}
	for _, topic := range orderTopics {
		bus.Subscribe(ctx, topic, func(event events.DomainEvent) error {
			orderID, err := types.ParseOrderID(event.AggregateID())
			if err != nil {
				log.Printf("[subscription] invalid order ID in event: %v", err)
				return nil
			}
			order, err := orderService.GetOrder(ctx, orderapp.GetOrderQuery{OrderID: orderID})
			if err != nil {
				log.Printf("[subscription] failed to fetch order %s: %v", event.AggregateID(), err)
				return nil
			}
			subMgr.Notify(event.AggregateID(), gql.ToGQLOrder(order))
			return nil
		})
	}

	// GraphQLのセットアップ
	resolver := &gql.Resolver{
		OrderService:      orderService,
		RestaurantService: restaurantService,
		SubscriptionMgr:   subMgr,
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql", gql.NewHandler(resolver))
	mux.Handle("/", gql.NewPlaygroundHandler("/graphql"))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Println("[api] listening on http://localhost:8080")
		log.Println("[api] GraphQL Playground: http://localhost:8080/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[api] server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[api] shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[api] shutdown error: %v", err)
	}
}
