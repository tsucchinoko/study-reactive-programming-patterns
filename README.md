# Food Delivery Order Tracking System

リアクティブプログラミングの学習を目的としたフードデリバリー注文追跡システムです。
GraphQL Subscription、RxJS、SQS、Redis Pub/Sub などのリアクティブパターンを実践的に学びます。

> 全てのデータはシミュレーションであり、実際の注文・配達は発生しません。

## 技術スタック

| レイヤー | 技術 |
|---|---|
| バックエンド | Go (gqlgen, pgx) |
| フロントエンド | React Native (Expo 54) + RxJS 7 + Apollo Client 4 |
| データベース | PostgreSQL 16 |
| メッセージング | Redis Pub/Sub, SQS (LocalStack) |
| インフラ | Docker Compose |
| モノレポ管理 | pnpm workspaces |

## プロジェクト構成

```
packages/
├── backend/           # Go API サーバー + シミュレーター
├── mobile/            # React Native (Expo) アプリ
└── infrastructure/    # SQL マイグレーション, Docker 初期化スクリプト
```

## セットアップ

### 前提条件

- Go 1.24+
- Node.js 20+ / pnpm 10+
- Docker / Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) (DBマイグレーション用)
- Expo Go アプリ (モバイル実機確認用、任意)

### 1. インフラ起動

```bash
docker compose up -d
```

PostgreSQL (5432)、Redis (6379)、LocalStack/SQS (4566) が起動します。

### 2. DBマイグレーション

```bash
cd packages/backend
make migrate-up
```

### 3. APIサーバー起動

```bash
cd packages/backend
make api
```

`http://localhost:8080` で GraphQL Playground にアクセスできます。

### 4. シミュレーター起動

```bash
cd packages/backend
make simulator
```

20件の日本語レストランデータをシードし、10秒間隔で注文を自動生成、3秒間隔で注文状態を遷移させます。

### 5. モバイルアプリ起動

```bash
cd packages/mobile
pnpm install
pnpm start
```

## 設計方針

### DDD (ドメイン駆動設計)

各境界づけられたコンテキスト (`order/`, `restaurant/`, `delivery/`) が独立したパッケージとして構成され、それぞれ `domain/`, `application/`, `infrastructure/` の3層を持ちます。

### 関数型プログラミング

- `Result[T]`, `Option[T]` モナドによるエラーハンドリング
- `fp.Map/Filter/Reduce/Pipe` による宣言的なデータ変換
- ドメインオブジェクトは不変 (状態変更メソッドは新しいインスタンスを返す)

### リアクティブデータフロー

```
OrderService
  → DomainEvent を InMemoryBus に発行
    → SubscriptionManager に通知
      → GraphQL WebSocket でクライアントに配信
        → Apollo Subscription → RxJS Observable → React State
```

## 注文ステートマシン

```
CREATED → CONFIRMED → PREPARING → READY → PICKED_UP → DELIVERING → DELIVERED
                                     ↑
                            PICKED_UP より前でキャンセル可能
```

## GraphQL API

### Query
- `order(id)` — 注文取得
- `orders(limit, offset)` — 注文一覧
- `ordersByStatus(status)` — ステータス別注文一覧
- `restaurant(id)` — レストラン取得
- `restaurants` — レストラン一覧

### Mutation
- `placeOrder(input)` — 注文作成
- `confirmOrder(orderId)` — 注文確認
- `cancelOrder(orderId, reason)` — 注文キャンセル
- `transitionOrder(orderId, newStatus)` — 注文状態遷移

### Subscription
- `orderStatusChanged(orderId)` — 注文状態のリアルタイム更新

## 開発コマンド

### バックエンド (`packages/backend/`)

| コマンド | 説明 |
|---|---|
| `make api` | APIサーバー起動 |
| `make simulator` | シミュレーター起動 |
| `make test` | テスト実行 |
| `make lint` | 静的解析 |
| `make generate` | GraphQL コード生成 (スキーマ変更後に実行) |
| `make migrate-up` | マイグレーション適用 |
| `make migrate-down` | マイグレーションロールバック |

### モバイル (`packages/mobile/`)

| コマンド | 説明 |
|---|---|
| `pnpm start` | Expo 開発サーバー起動 |
| `pnpm ios` | iOS シミュレーター起動 |
| `pnpm android` | Android エミュレーター起動 |

## フェーズ計画

| フェーズ | 内容 | 状態 |
|---|---|---|
| Phase 1 | Go バックエンド + DDD + GraphQL API + シミュレーター | 完了 |
| Phase 2 | React Native + GraphQL Subscription + RxJS | 完了 |
| Phase 3 | Redis Pub/Sub + SQS + 配達コンテキスト | 未着手 |
| Phase 4 | アナリティクス + カオスエンジニアリング + オブザーバビリティ | 未着手 |
