package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	deliveryapp "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/application"
	deliverydomain "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	deliveryinfra "github.com/tsucchinoko/food-delivery-tracker/internal/delivery/infrastructure"
	gql "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/graphql"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/postgres"
	infraredis "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/redis"
	infrasqs "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/sqs"
	orderapp "github.com/tsucchinoko/food-delivery-tracker/internal/order/application"
	orderdomain "github.com/tsucchinoko/food-delivery-tracker/internal/order/domain"
	orderinfra "github.com/tsucchinoko/food-delivery-tracker/internal/order/infrastructure"
	restaurantapp "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/application"
	restaurantinfra "github.com/tsucchinoko/food-delivery-tracker/internal/restaurant/infrastructure"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[api] starting...")

	// イベント型をレジストリに登録（デシリアライズに必要）
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
		log.Fatalf("[api] failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Redis クライアントのセットアップ
	redisClient := infraredis.NewClient(infraredis.DefaultConfig())
	defer redisClient.Close()

	// SQS クライアントのセットアップ
	sqsClient := infrasqs.NewClient(infrasqs.DefaultEndpoint(), infrasqs.DefaultRegion())
	orderQueueURL, err := infrasqs.ResolveQueueURL(ctx, sqsClient, "order-events-queue")
	if err != nil {
		log.Fatalf("[api] failed to resolve SQS queue URL: %v", err)
	}

	// イベントパブリッシャーの合成: Redis Pub/Sub + SQS
	publisher := events.ComposePublishers(
		infraredis.NewPublisher(redisClient),
		infrasqs.NewPublisher(sqsClient, orderQueueURL),
	)

	// GraphQL Subscription用のサブスクリプションマネージャーをセットアップ
	subMgr := gql.NewSubscriptionManager()

	// リポジトリのセットアップ
	orderRepo := orderinfra.NewPostgresOrderRepository(pool)
	restaurantRepo := restaurantinfra.NewPostgresRestaurantRepository(pool)
	driverRepo := deliveryinfra.NewPostgresDriverRepository(pool)
	assignmentRepo := deliveryinfra.NewPostgresAssignmentRepository(pool)

	// サービスのセットアップ
	orderService := orderapp.NewOrderService(orderRepo, publisher)
	restaurantService := restaurantapp.NewRestaurantService(restaurantRepo)
	deliveryService := deliveryapp.NewDeliveryService(driverRepo, assignmentRepo, publisher)

	// Redis Subscriber → サブスクリプションマネージャーの接続
	// 注文イベント発火時に注文を再取得し、サブスクライバーに通知する
	redisSub := infraredis.NewSubscriber(redisClient)
	orderTopics := []string{"order.placed", "order.confirmed", "order.status_changed", "order.cancelled"}
	for _, topic := range orderTopics {
		topic := topic // ループ変数キャプチャ
		go func() {
			err := redisSub.Subscribe(ctx, topic, func(event events.DomainEvent) error {
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
			if err != nil && ctx.Err() == nil {
				log.Printf("[api] redis subscriber for %s stopped: %v", topic, err)
			}
		}()
	}

	// Redis PSubscribe でドライバー位置更新をパターン購読
	// driver.location.{driverID} チャンネルの全メッセージを受信し、SubscriptionManager 経由で WebSocket に配信する
	go func() {
		pubsub := redisClient.PSubscribe(ctx, "driver.location.*")
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				// チャンネル名から driverID を抽出: "driver.location.{driverID}"
				parts := strings.SplitN(msg.Channel, ".", 3)
				if len(parts) < 3 {
					continue
				}
				driverID := parts[2]

				// イベントエンベロープからペイロードを抽出
				var envelope struct {
					Payload json.RawMessage `json:"payload"`
				}
				if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
					log.Printf("[api] failed to parse location envelope: %v", err)
					continue
				}
				var loc struct {
					DriverID  string  `json:"driver_id"`
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				}
				if err := json.Unmarshal(envelope.Payload, &loc); err != nil {
					log.Printf("[api] failed to parse location payload: %v", err)
					continue
				}

				gqlLoc := gql.ToGQLDriverLocation(driverID, loc.Latitude, loc.Longitude, time.Now())
				subMgr.NotifyDriverLocation(driverID, gqlLoc)
			}
		}
	}()

	// GraphQLのセットアップ
	resolver := &gql.Resolver{
		OrderService:      orderService,
		RestaurantService: restaurantService,
		DeliveryService:   deliveryService,
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
