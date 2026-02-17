.PHONY: api worker simulator test lint generate migrate-up migrate-down

api:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

simulator:
	go run cmd/simulator/main.go

test:
	go test ./... -v

lint:
	go vet ./...

generate:
	go run github.com/99designs/gqlgen generate

migrate-up:
	migrate -path internal/infrastructure/postgres/migrations -database "postgres://delivery:delivery@localhost:5432/food_delivery?sslmode=disable" up

migrate-down:
	migrate -path internal/infrastructure/postgres/migrations -database "postgres://delivery:delivery@localhost:5432/food_delivery?sslmode=disable" down
