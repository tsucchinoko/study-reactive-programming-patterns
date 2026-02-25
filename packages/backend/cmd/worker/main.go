package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	infrasqs "github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/sqs"
	notificationapp "github.com/tsucchinoko/food-delivery-tracker/internal/notification/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[worker] starting SQS consumer worker...")

	// SQS クライアントのセットアップ
	sqsClient := infrasqs.NewClient(infrasqs.DefaultEndpoint(), infrasqs.DefaultRegion())

	// order-events-queue の URL を取得
	orderQueueURL, err := infrasqs.ResolveQueueURL(ctx, sqsClient, "order-events-queue")
	if err != nil {
		log.Fatalf("[worker] failed to resolve order-events-queue URL: %v", err)
	}

	// notification-events-queue の URL を取得
	notifQueueURL, err := infrasqs.ResolveQueueURL(ctx, sqsClient, "notification-events-queue")
	if err != nil {
		log.Fatalf("[worker] failed to resolve notification-events-queue URL: %v", err)
	}

	// 注文イベントキューの Consumer（通知処理）
	orderConsumer := infrasqs.NewConsumer(sqsClient, orderQueueURL, notificationapp.OrderEventHandler)

	// 通知イベントキューの Consumer（同じく通知処理）
	notifConsumer := infrasqs.NewConsumer(sqsClient, notifQueueURL, notificationapp.OrderEventHandler)

	// 複数の Consumer を並行起動
	go func() {
		log.Printf("[worker] starting order-events consumer (queue: %s)", orderQueueURL)
		if err := orderConsumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("[worker] order-events consumer stopped: %v", err)
		}
	}()

	go func() {
		log.Printf("[worker] starting notification-events consumer (queue: %s)", notifQueueURL)
		if err := notifConsumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("[worker] notification-events consumer stopped: %v", err)
		}
	}()

	log.Println("[worker] running — press Ctrl+C to stop")

	<-ctx.Done()
	log.Println("[worker] shutting down...")
}
