package sqs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
)

// NewPublisher は DomainEvent を SQS キューに送信する EventPublisher を返す。
func NewPublisher(client *sqs.Client, queueURL string) events.EventPublisher {
	return func(ctx context.Context, event events.DomainEvent) error {
		data, err := events.SerializeEvent(event)
		if err != nil {
			return fmt.Errorf("sqs publish: serialize: %w", err)
		}

		_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(string(data)),
			MessageAttributes: map[string]sqstypes.MessageAttributeValue{
				"EventType": {
					DataType:    aws.String("String"),
					StringValue: aws.String(event.EventType()),
				},
				"Topic": {
					DataType:    aws.String("String"),
					StringValue: aws.String(event.Topic()),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("sqs send to %s: %w", queueURL, err)
		}

		log.Printf("[sqs] sent event %s to queue %s", event.EventType(), queueURL)
		return nil
	}
}
