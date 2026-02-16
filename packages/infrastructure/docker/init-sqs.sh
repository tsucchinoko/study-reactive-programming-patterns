#!/bin/bash
awslocal sqs create-queue --queue-name order-events-queue
awslocal sqs create-queue --queue-name notification-events-queue
awslocal sqs create-queue --queue-name analytics-events-queue
awslocal sqs create-queue --queue-name order-events-dlq
echo "SQS queues created successfully"
