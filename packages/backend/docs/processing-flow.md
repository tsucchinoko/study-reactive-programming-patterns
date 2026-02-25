

バックエンドの処理フローを調査します。

## バックエンドの処理フロー

### 全体構成

3つのエントリポイントがあります:

| エントリポイント | 役割 |
|---|---|
| `cmd/api/main.go` | APIサーバー (port 8080) |
| `cmd/simulator/main.go` | 注文シミュレーター |
| `cmd/worker/main.go` | SQSワーカー (Phase 3用スタブ) |

---

### フローA: GraphQL Mutation（例: placeOrder）

```
クライアント HTTP POST /graphql
  → gqlgen が mutationResolver.PlaceOrder() にディスパッチ
    → 入力をパースし PlaceOrderCommand を構築
    → OrderService.PlaceOrder()
      → domain.BuildOrder()  ← バリデーション + ドメインオブジェクト生成
      → repo.Save()          ← PostgreSQL に INSERT (UPSERT)
      → bus.Publish(OrderPlaced イベント)
        → イベントハンドラ: DBから注文を再取得 → subMgr.Notify()
          → 該当orderIDの全サブスクライバーchannelにpush
    → ToGQLOrder() でドメイン→GraphQLモデルに変換
  → JSON レスポンス返却
```

### フローB: GraphQL Subscription（orderStatusChanged）

```
クライアント WebSocket /graphql (graphql-ws プロトコル)
  → subscriptionResolver.OrderStatusChanged()
    → subMgr.Subscribe(orderID, subscriberID) → channelを返す
    → goroutine: ctx.Done()待ち → クライアント切断時にunsubscribe()
    → channelをgqlgenに返す
  → gqlgen がchannelから読み取り → WebSocketメッセージとして送信
```

### フローC: シミュレーター

```
10秒ごと:
  GenerateRandomOrder() → ランダムにレストラン・メニュー選択 → PlaceOrder()

3秒ごと:
  RunStateTransitions() → アクティブな全注文を走査:
    - 5%の確率: CancelOrder() (PICKED_UP前のみ)
    - それ以外: ステータスに応じた遅延後、次の状態に遷移
```

---

### 注文ステートマシン

```
CREATED → CONFIRMED → PREPARING → READY → PICKED_UP → DELIVERING → DELIVERED
  |          |           |          |
  +----------+-----------+----------+--→ CANCELLED
```

各遷移は `order_status.go` の `validTransitions` マップで厳密に制御されています。ドメインオブジェクトは**イミュータブル**で、`Confirm()`/`TransitionTo()`/`Cancel()` は常に**新しいインスタンス**を返します。

---

### レイヤー構成（DDD + パッケージバイフィーチャー）

```
order/
  domain/        ← 純粋なドメインロジック（副作用なし）
    order.go          状態遷移メソッド（新インスタンスを返す）
    order_status.go   ステートマシン定義
    events.go         ドメインイベント定義
    service.go        BuildOrder（バリデーション+生成）
  application/   ← ユースケースのオーケストレーション
    service.go        リポジトリ保存 + イベント発行の調整
  infrastructure/← DB永続化
    postgres_repository.go  UPSERT + N+1クエリパターン
```

---

### イベントバスの仕組み

`InMemoryBus` がトピックベースのPub/Subを提供:

- **4トピック**を購読: `order.placed`, `order.confirmed`, `order.status_changed`, `order.cancelled`
- 各ハンドラ: イベント受信 → DBから最新注文を再取得 → `SubscriptionManager.Notify()` でWebSocketにpush
- `Notify()` はバッファサイズ1のchannelに**非ブロッキング送信**（遅いクライアントがイベント処理をブロックしない設計）

---

### 関数型プログラミングの要素

| パッケージ | 役割 |
|---|---|
| `shared/result/` | `Result[T]` モナド (`Ok`/`Err`/`Map`/`FlatMap`) |
| `shared/option/` | `Option[T]` モナド (`Some`/`None`) |
| `shared/fp/` | `Map`/`Filter`/`Reduce`/`Pipe` ユーティリティ |
| `shared/types/` | 型付きUUID (`OrderID`等), `Money`(JPY), `Timestamp` |

ドメイン層のバリデーションチェーンは `Result.FlatMap` で関数合成されており、Go言語でありながら関数型のスタイルを採用しています。

---

### 注意点

シミュレーターとAPIサーバーは**別プロセス**で動作し、それぞれ独立した`InMemoryBus`を持っています。そのため、シミュレーターが発行したイベントはAPIサーバーのSubscriptionManagerには届きません。Phase 3でRedis Pub/Subに置き換えることで、プロセス間のイベント共有が実現される予定です。
