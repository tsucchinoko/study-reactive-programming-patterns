# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Food delivery order tracking system for learning reactive programming patterns (GraphQL Subscriptions, RxJS, SQS, Redis Pub/Sub). All data is simulated — no real orders or deliveries.

**Monorepo** managed by pnpm workspaces (`packages/*`):
- `packages/backend/` — Go API server + simulator
- `packages/mobile/` — React Native (Expo) app
- `packages/infrastructure/` — SQL migrations, Docker init scripts

## Commands

### Infrastructure
```bash
docker compose up -d          # Start PostgreSQL, Redis, LocalStack (SQS)
```

### Backend (run from `packages/backend/`)
```bash
make api                      # Run API server (port 8080)
make simulator                # Run order simulator (seeds restaurants + generates orders)
make test                     # go test ./... -v
make lint                     # go vet ./...
make generate                 # Regenerate gqlgen code from schema.graphql
make migrate-up               # Apply DB migrations
make migrate-down             # Rollback DB migrations
```

### Mobile (run from `packages/mobile/`)
```bash
pnpm start                    # Expo dev server
pnpm ios                      # iOS simulator
pnpm android                  # Android emulator
```

## Architecture

### Design Principles
- **DDD with package-by-feature**: each bounded context (`order/`, `restaurant/`, `delivery/`) contains `domain/`, `application/`, `infrastructure/` layers
- **Immutability**: domain struct fields are unexported; state-change methods (`Confirm`, `Cancel`, `TransitionTo`) return new instances
- **Functional programming**: `Result[T]` and `Option[T]` monads in `internal/shared/`; `fp.Map/Filter/Reduce/Pipe` utilities throughout
- **Pure domain logic**: no side effects in domain layer; all IO in application/infrastructure layers

### Backend (Go)

**Entry points:**
- `cmd/api/main.go` — wires repositories, services, event bus, SubscriptionManager, and GraphQL handler
- `cmd/simulator/main.go` — seeds 20 Japanese restaurants, generates random orders every 10s, advances order states every 3s
- `cmd/worker/main.go` — stub for Phase 3 SQS consumer

**Event flow:**
`OrderService` → publishes `DomainEvent` to `InMemoryBus` → bus handler fetches updated order → `SubscriptionManager.Notify()` → pushes to subscriber channels → GraphQL WebSocket transport delivers to client

**Key wiring in `cmd/api/main.go`:** The event bus subscribes to 4 order topics (`order.placed`, `order.confirmed`, `order.status_changed`, `order.cancelled`). Each handler fetches the full order and notifies the SubscriptionManager.

**GraphQL:**
- Schema: `internal/infrastructure/graphql/schema.graphql`
- gqlgen config: `gqlgen.yml`
- Generated code: `internal/infrastructure/graphql/generated/` (do not edit manually)
- Resolvers: `internal/infrastructure/graphql/schema.resolvers.go`
- After schema changes, run `make generate` to regenerate

**Order state machine** (`internal/order/domain/order_status.go`): `CREATED → CONFIRMED → PREPARING → READY → PICKED_UP → DELIVERING → DELIVERED`. Cancellation allowed before `PICKED_UP`. The `validTransitions` map enforces legal transitions.

**Shared Kernel** (`internal/shared/`):
- `result/` — `Result[T]` monad (`Ok`, `Err`, `Map`, `FlatMap`, `Match`)
- `option/` — `Option[T]` monad (`Some`, `None`, `Map`, `FlatMap`, `Filter`)
- `fp/` — generic functional utilities
- `types/` — typed UUIDs (`OrderID`, `CustomerID`, etc.), `Money` (JPY), `Timestamp`
- `events/` — `DomainEvent` interface, `InMemoryBus` (to be replaced by Redis Pub/Sub in Phase 3)

### Mobile (React Native / Expo)

**Data flow:** Apollo Client split link — HTTP for queries/mutations, WebSocket (`graphql-ws`) for subscriptions. RxJS bridges Apollo subscriptions to React state.

**Key patterns:**
- `shared/rx/operators.ts` — `fromApolloSubscription` (Apollo Observable → RxJS Observable), `retryWithBackoff` (exponential backoff)
- `shared/rx/useObservable.ts` — `useObservable<T>` hook: subscribes to Observable, manages `{status, data, error}` state, auto-unsubscribes on unmount
- `features/order/hooks/useOrderStream.ts` — composes subscription + RxJS operators for live order updates
- Feature structure: `features/<context>/graphql/`, `hooks/`, `components/`, `screens/`

**Import conventions for Apollo Client v4:**
- `useQuery`, `useMutation`, `ApolloProvider` from `@apollo/client/react`
- `gql`, `ApolloClient`, link utilities from `@apollo/client`

## Database

PostgreSQL 16, credentials: `delivery/delivery`, database: `food_delivery`, port 5432.
Migrations in `packages/backend/internal/infrastructure/postgres/migrations/` using golang-migrate.

## Phase Roadmap

- **Phase 1** (done): Go backend + DDD + GraphQL API + Simulator
- **Phase 2** (done): React Native + GraphQL Subscription + RxJS
- **Phase 3** (next): Redis Pub/Sub replaces InMemoryBus, SQS for async processing, Delivery bounded context
- **Phase 4**: Analytics, Chaos engineering, Observability
