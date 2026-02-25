package sqs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// MessageHandler は SQS メッセージを処理する関数型。
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer は SQS キューからメッセージをポーリングして処理する。
type Consumer struct {
	client   *sqs.Client
	queueURL string
	handler  MessageHandler
}

// NewConsumer は新しい SQS Consumer を生成する。
func NewConsumer(client *sqs.Client, queueURL string, handler MessageHandler) *Consumer {
	return &Consumer{
		client:   client,
		queueURL: queueURL,
		handler:  handler,
	}
}

// Run はメッセージのポーリングループを開始する。
// コンテキストがキャンセルされるまでブロッキングで動作する。
func (c *Consumer) Run(ctx context.Context) error {
	log.Printf("[sqs] consumer started for queue %s", c.queueURL)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[sqs] consumer stopping for queue %s", c.queueURL)
			return ctx.Err()
		default:
			if err := c.poll(ctx); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				log.Printf("[sqs] poll error: %v, retrying in 5s", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (c *Consumer) poll(ctx context.Context) error {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20, // ロングポーリング
		MessageAttributeNames: []string{
			"EventType",
			"Topic",
		},
	})
	if err != nil {
		return fmt.Errorf("receive messages: %w", err)
	}

	for _, msg := range out.Messages {
		if err := c.processMessage(ctx, msg); err != nil {
			log.Printf("[sqs] failed to process message %s: %v", *msg.MessageId, err)
			continue
		}

		// 処理成功: メッセージを削除
		_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(c.queueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("[sqs] failed to delete message %s: %v", *msg.MessageId, err)
		}
	}

	return nil
}

func (c *Consumer) processMessage(ctx context.Context, msg sqstypes.Message) error {
	if msg.Body == nil {
		return fmt.Errorf("message body is nil")
	}
	return c.handler(ctx, []byte(*msg.Body))
}
