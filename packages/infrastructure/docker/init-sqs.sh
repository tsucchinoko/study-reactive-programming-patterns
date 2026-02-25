#!/bin/bash
# LocalStack の新バージョンでは DEFAULT_REGION が無視されるため、明示的にリージョンを指定する
awslocal sqs create-queue --queue-name order-events-queue --region ap-northeast-1
awslocal sqs create-queue --queue-name notification-events-queue --region ap-northeast-1
awslocal sqs create-queue --queue-name analytics-events-queue --region ap-northeast-1
awslocal sqs create-queue --queue-name order-events-dlq --region ap-northeast-1
echo "SQS queues created successfully"
