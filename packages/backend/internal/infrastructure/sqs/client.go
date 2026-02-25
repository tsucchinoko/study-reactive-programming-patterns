package sqs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// DefaultEndpoint は環境変数 SQS_ENDPOINT またはデフォルト LocalStack エンドポイントを返す。
func DefaultEndpoint() string {
	endpoint := os.Getenv("SQS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}
	return endpoint
}

// DefaultRegion は環境変数 AWS_REGION またはデフォルト ap-northeast-1 を返す。
func DefaultRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	return region
}

// NewClient は LocalStack 対応の SQS クライアントを生成する。
func NewClient(endpoint, region string) *sqs.Client {
	return sqs.New(sqs.Options{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider("test", "test", ""),
	})
}

// ResolveQueueURL はキュー名から URL を取得する。
func ResolveQueueURL(ctx context.Context, client *sqs.Client, queueName string) (string, error) {
	out, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", fmt.Errorf("resolve queue URL for %s: %w", queueName, err)
	}
	return *out.QueueUrl, nil
}
