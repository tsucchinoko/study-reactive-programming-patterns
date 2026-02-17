# Food Delivery Order Tracking System - 実装計画

## Context

リアクティブプログラミング（GraphQL Subscription, RxJS, SQS, Redis Pub/Sub）を実践的に学ぶための学習プロジェクト。フードデリバリーの注文追跡システムをシミュレーションベースで構築する。

**設計原則**: DDD, Package by Feature, 関数型プログラミング  
**技術スタック**: Go (Backend) + React Native/Expo (Frontend) + PostgreSQL + Redis + SQS (LocalStack)

---

## アーキテクチャ概要

### Bounded Contexts (DDD)

| Context | 責務 | 主要技術 |
|---|---|---|
| **Order** | 注文ライフサイクル管理 | GraphQL Mutation/Subscription, Redis Pub/Sub, SQS |
| **Restaurant** | 店舗・メニュー・調理シミュレーション | GraphQL Query |
| **Delivery** | 配達員アサイン・位置追跡シミュレーション | GraphQL Subscription, Redis Pub/Sub |
| **Notification** | 通知の非同期配信 | SQS Consumer |
| **Analytics** | リアルタイム統計・ダッシュボード | SQS Consumer, GraphQL Subscription, RxJS |

### 注文状態遷移

```
Created → Confirmed → Preparing → Ready → PickedUp → Delivering → Delivered
                                                                    ↗
                        (任意のタイミングで) → Cancelled
```

---

## モノレポ構成

```
reactive-programming/
├── packages/
│   ├── backend/           # Go — API, Worker, Simulator
│   ├── mobile/            # React Native (Expo)
│   └── infrastructure/    # Docker Compose, LocalStack設定, マイグレーション
├── pnpm-workspace.yml
├── package.json
└── docker-compose.yml
```

> シミュレータはGoバックエンド内 (`cmd/simulator`) に統一する。TypeScript側には別途シミュレータを置かない。

---

## Backend ディレクトリ構造 (Go)

```
packages/backend/
├── cmd/
│   ├── api/main.go              # GraphQL APIサーバー
│   ├── worker/main.go           # SQS Consumerワーカー
│   └── simulator/main.go        # 注文・配達シミュレータ
├── internal/
│   ├── order/                   # === Order Bounded Context ===
│   │   ├── domain/
│   │   │   ├── order.go              # Aggregate Root (不変)
│   │   │   ├── order_item.go         # Value Object
│   │   │   ├── order_status.go       # Value Object (enum + 状態遷移ルール)
│   │   │   ├── events.go             # Domain Events
│   │   │   ├── repository.go         # Repository Interface
│   │   │   └── service.go            # Pure Functions (計算・検証)
│   │   ├── application/
│   │   │   ├── commands.go           # PlaceOrder, ConfirmOrder, CancelOrder
│   │   │   ├── queries.go            # GetOrder, ListOrders
│   │   │   └── service.go            # Application Service (副作用のオーケストレーション)
│   │   ├── infrastructure/
│   │   │   ├── postgres_repository.go
│   │   │   ├── redis_publisher.go
│   │   │   └── sqs_producer.go
│   │   └── graphql/
│   │       ├── resolvers.go
│   │       └── subscriptions.go
│   ├── restaurant/              # === Restaurant Bounded Context ===
│   │   ├── domain/
│   │   │   ├── restaurant.go         # Aggregate Root
│   │   │   ├── menu.go               # Menu Aggregate
│   │   │   ├── menu_item.go          # Value Object
│   │   │   ├── events.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── application/
│   │   │   ├── commands.go
│   │   │   ├── queries.go
│   │   │   └── service.go
│   │   ├── infrastructure/
│   │   │   └── postgres_repository.go
│   │   └── graphql/
│   │       └── resolvers.go
│   ├── delivery/                # === Delivery Bounded Context ===
│   │   ├── domain/
│   │   │   ├── driver.go             # Aggregate Root
│   │   │   ├── assignment.go         # Entity
│   │   │   ├── location.go           # Value Object
│   │   │   ├── events.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── application/
│   │   │   ├── commands.go
│   │   │   ├── queries.go
│   │   │   └── service.go
│   │   ├── infrastructure/
│   │   │   ├── postgres_repository.go
│   │   │   ├── redis_publisher.go
│   │   │   └── location_simulator.go  # Random Walk
│   │   └── graphql/
│   │       ├── resolvers.go
│   │       └── subscriptions.go
│   ├── notification/            # === Notification Bounded Context ===
│   │   ├── domain/
│   │   │   ├── notification.go
│   │   │   ├── events.go
│   │   │   └── service.go
│   │   ├── application/
│   │   │   ├── handlers.go           # SQSイベントハンドラ
│   │   │   └── service.go
│   │   └── infrastructure/
│   │       ├── sqs_consumer.go
│   │       └── postgres_repository.go
│   ├── analytics/               # === Analytics Bounded Context ===
│   │   ├── domain/
│   │   │   ├── metric.go
│   │   │   ├── aggregation.go        # Pure: 集計関数
│   │   │   └── service.go
│   │   ├── application/
│   │   │   ├── handlers.go
│   │   │   ├── queries.go
│   │   │   └── service.go
│   │   ├── infrastructure/
│   │   │   ├── sqs_consumer.go
│   │   │   ├── redis_cache.go
│   │   │   └── postgres_repository.go
│   │   └── graphql/
│   │       ├── resolvers.go
│   │       └── subscriptions.go
│   ├── shared/                  # === Shared Kernel ===
│   │   ├── result/result.go          # Result[T] モナド
│   │   ├── option/option.go          # Option[T] モナド
│   │   ├── fp/pipe.go                # Pipe, Compose, Map, Filter, Reduce
│   │   ├── events/
│   │   │   ├── event.go              # DomainEvent interface
│   │   │   ├── bus.go                # EventBus interface
│   │   │   └── dispatcher.go         # Redis + SQS へのディスパッチ
│   │   └── types/
│   │       ├── id.go                 # Typed ID (OrderID, DriverID等)
│   │       ├── money.go              # Money Value Object
│   │       └── timestamp.go
│   └── infrastructure/          # === Cross-cutting ===
│       ├── graphql/
│       │   ├── server.go
│       │   ├── schema.graphql
│       │   └── generated/            # gqlgen生成コード
│       ├── postgres/
│       │   ├── client.go
│       │   └── transaction.go        # 関数型トランザクションラッパー
│       ├── redis/
│       │   ├── client.go
│       │   └── pubsub.go
│       └── sqs/
│           ├── client.go
│           ├── producer.go
│           └── consumer.go
├── go.mod
├── go.sum
├── gqlgen.yml
└── Makefile
```

---

## Frontend ディレクトリ構造 (React Native / Expo)

```
packages/mobile/
├── src/
│   ├── features/                # === Package by Feature ===
│   │   ├── order/
│   │   │   ├── screens/
│   │   │   │   ├── OrderListScreen.tsx
│   │   │   │   ├── OrderDetailScreen.tsx
│   │   │   │   └── PlaceOrderScreen.tsx
│   │   │   ├── components/
│   │   │   │   ├── OrderCard.tsx
│   │   │   │   ├── OrderTimeline.tsx
│   │   │   │   └── StatusBadge.tsx
│   │   │   ├── hooks/
│   │   │   │   ├── useOrderList.ts
│   │   │   │   ├── useOrderDetail.ts
│   │   │   │   └── useOrderSubscription.ts
│   │   │   ├── services/
│   │   │   │   ├── orderService.ts       # GraphQL queries/mutations
│   │   │   │   └── orderStream.ts        # RxJS stream composition
│   │   │   └── types/
│   │   │       └── order.types.ts
│   │   ├── restaurant/
│   │   │   ├── screens/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── services/
│   │   │   └── types/
│   │   ├── delivery/
│   │   │   ├── screens/
│   │   │   │   └── TrackingScreen.tsx
│   │   │   ├── components/
│   │   │   │   ├── MapView.tsx
│   │   │   │   └── DriverMarker.tsx
│   │   │   ├── hooks/
│   │   │   │   └── useDriverLocation.ts  # RxJS location stream
│   │   │   ├── services/
│   │   │   │   └── locationStream.ts     # 補間・スロットリング
│   │   │   └── types/
│   │   └── analytics/
│   │       ├── screens/
│   │       │   └── DashboardScreen.tsx
│   │       ├── components/
│   │       │   ├── MetricCard.tsx
│   │       │   └── RealtimeChart.tsx
│   │       ├── hooks/
│   │       │   └── useDashboardMetrics.ts  # combineLatest + scan
│   │       └── services/
│   │           └── metricsStream.ts
│   ├── shared/
│   │   ├── graphql/
│   │   │   ├── client.ts                # Apollo Client (WebSocket link)
│   │   │   ├── queries/                 # .graphql ファイル
│   │   │   └── subscriptions/           # .graphql ファイル
│   │   ├── rx/
│   │   │   ├── operators/
│   │   │   │   ├── retryWithBackoff.ts
│   │   │   │   └── bufferWithBackpressure.ts
│   │   │   └── hooks/
│   │   │       ├── useObservable.ts
│   │   │       └── useSubscription.ts
│   │   └── components/
│   │       ├── ErrorBoundary.tsx
│   │       └── LoadingSpinner.tsx
│   ├── navigation/
│   │   └── RootNavigator.tsx
│   └── App.tsx
├── app.json
├── package.json
└── tsconfig.json
```

---

## 関数型プログラミング: Go での主要パターン

### Result[T] モナド
```go
type Result[T any] struct { value T; err error }

func (r Result[T]) Map(f func(T) T) Result[T]
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T]
func (r Result[T]) Match(onOk func(T), onErr func(error))
```

### 不変ドメインオブジェクト
```go
// With メソッドで新しいコピーを返す (mutation なし)
func (o Order) WithStatus(s OrderStatus) Order { ... }
```

### 純粋関数によるドメインロジック
```go
// 副作用なし — テスト容易
func CanTransitionTo(current, next OrderStatus) Result[Unit]
func CalculateOrderTotal(items []OrderItem) Money
func SelectBestDriver(drivers []Driver, pickup Location) Option[Driver]
```

### 関数合成によるイベント配信
```go
type EventPublisher func(ctx context.Context, event DomainEvent) Result[Unit]

func ComposePublishers(pubs ...EventPublisher) EventPublisher { ... }

// 使用例
publisher := ComposePublishers(PublishToRedis(redis), PublishToSQS(sqs, url))
```

---

## RxJS: Frontend での主要パターン

| パターン | 使用箇所 | 主要オペレータ |
|---|---|---|
| ストリーム購読 + 自動再接続 | 注文ステータス | retryWhen, delay, shareReplay |
| スロットリング + 補間 | 配達員位置追跡 | throttleTime, scan, switchMap, interval |
| 複数ストリーム合成 | ダッシュボード | combineLatest, bufferTime, distinctUntilChanged |
| 移動平均 | チャート | scan (window), map |
| バックプレッシャー | 高頻度更新 | bufferTime + filter |

---

## フェーズ別実装計画

### Phase 1: コアドメイン + Goバックエンド基盤

**ゴール**: DDD構造の確立、Order/Restaurantコンテキスト、GraphQL API、注文シミュレータ

**作業内容**:
1. プロジェクトセットアップ (Go modules, Docker Compose for PostgreSQL)
2. Shared Kernel 実装 (Result, Option, FP utilities, Typed IDs, Money)
3. Order Context: ドメインモデル, 純粋関数, PostgreSQLリポジトリ, GraphQL resolvers
4. Restaurant Context: ドメインモデル, リポジトリ, シードデータ(20店舗)
5. GraphQL スキーマ + gqlgen セットアップ (Query, Mutation)
6. 注文シミュレータ (`cmd/simulator`): ランダム注文生成 + 状態遷移自動進行

**検証**: シミュレータ起動 → 注文が自動生成 → GraphQL queryで確認可能

### Phase 2: モバイルアプリ + リアルタイムSubscription

**ゴール**: React Native app, GraphQL Subscription, RxJSストリーム処理

**作業内容**:
1. Expo プロジェクト初期化 + Apollo Client (WebSocket link)
2. GraphQL Subscription 実装 (Backend: in-memory pub/sub → Phase 3でRedisに移行)
3. Order Feature: 一覧/詳細画面, useOrderSubscription hook, orderStream (RxJS)
4. Restaurant Feature: 一覧/メニュー画面
5. RxJS utilities: retryWithBackoff, useObservable hook
6. Navigation セットアップ

**検証**: アプリで注文一覧表示 → シミュレータが状態遷移 → リアルタイムでUI更新

### Phase 3: Redis Pub/Sub + SQS + Delivery

**ゴール**: イベント駆動アーキテクチャ、非同期処理、配達追跡

**作業内容**:
1. Docker Compose に Redis + LocalStack (SQS) 追加
2. Event Infrastructure: Redis publisher, SQS producer, ComposePublishers
3. Delivery Context: ドメインモデル, Location Simulator (random walk), GraphQL Subscription
4. Order/Restaurant Context: 状態遷移時に Redis + SQS へイベント発行
5. Notification Context: SQS Consumer, モック通知
6. Mobile Delivery Feature: TrackingScreen, useDriverLocation (throttle + 補間)
7. SQS Worker (`cmd/worker`)

**検証**: 注文イベント → Redis + SQS に配信 → Worker処理 → 配達員位置がマップ上でスムーズに移動

### Phase 4: Analytics + Chaos + 可観測性

**ゴール**: リアルタイムダッシュボード、カオスエンジニアリング、可観測性

**作業内容**:
1. Analytics Context: SQS Consumer, Windowed集計, Redis Cache, GraphQL Subscription
2. Mobile Dashboard: MetricCard, RealtimeChart, combineLatest + scan
3. Chaos Engine: ランダム障害注入 (遅延, キャンセル, 配達失敗)
4. Observability: 構造化ログ (zerolog), OpenTelemetry tracing
5. Advanced RxJS: バックプレッシャー, メモリリーク防止, 移動平均

**検証**: ダッシュボードにリアルタイムメトリクス表示 → Chaos有効化 → アプリが障害を適切にハンドリング

---

## インフラ構成

### docker-compose.yml
- **PostgreSQL**: 5432
- **Redis**: 6379
- **LocalStack** (SQS): 4566

### SQS キュー (LocalStack)
- `order-events-queue`
- `notification-events-queue`
- `analytics-events-queue`
- `order-events-dlq` (Dead Letter Queue)

### Redis Pub/Sub チャンネル
- `order.events.{orderId}`
- `driver.location.{driverId}`
- `analytics.metrics`

---

## 主要ライブラリ

### Backend (Go)
- `github.com/99designs/gqlgen` — GraphQL
- `github.com/lib/pq` or `github.com/jackc/pgx/v5` — PostgreSQL
- `github.com/redis/go-redis/v9` — Redis
- `github.com/aws/aws-sdk-go-v2` — SQS (LocalStack)
- `github.com/golang-migrate/migrate/v4` — マイグレーション
- `github.com/rs/zerolog` — 構造化ログ

### Frontend (React Native / Expo)
- `rxjs` — リアクティブストリーム
- `@apollo/client` — GraphQL + Subscription
- `graphql-ws` — WebSocket transport
- `react-native-maps` — 地図 (Phase 3)

---

## 検証方法

各フェーズ完了時:
1. **Unit Tests**: 純粋関数のテスト (`go test ./internal/.../domain/...`)
2. **Integration Tests**: リポジトリテスト (testcontainers)
3. **E2E確認**: シミュレータ起動 → GraphQL Playground/アプリで動作確認
4. **RxJS Tests**: TestScheduler を使ったストリームテスト
