package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gql "github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/postgres"
	orderapp "github.com/daichitsuchiya/food-delivery-tracker/internal/order/application"
	orderinfra "github.com/daichitsuchiya/food-delivery-tracker/internal/order/infrastructure"
	restaurantapp "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/application"
	restaurantinfra "github.com/daichitsuchiya/food-delivery-tracker/internal/restaurant/infrastructure"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/shared/events"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("[api] starting...")

	// Connect to PostgreSQL
	pool, err := postgres.NewPool(ctx, postgres.DefaultConfig())
	if err != nil {
		log.Fatalf("[api] failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Set up event bus (in-memory for Phase 1)
	bus := events.NewInMemoryBus()

	// Set up repositories
	orderRepo := orderinfra.NewPostgresOrderRepository(pool)
	restaurantRepo := restaurantinfra.NewPostgresRestaurantRepository(pool)

	// Set up services
	orderService := orderapp.NewOrderService(orderRepo, bus.Publish)
	restaurantService := restaurantapp.NewRestaurantService(restaurantRepo)

	// Set up GraphQL
	resolver := &gql.Resolver{
		OrderService:      orderService,
		RestaurantService: restaurantService,
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql", gql.NewHandler(resolver))
	mux.Handle("/", gql.NewPlaygroundHandler("/graphql"))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Println("[api] listening on http://localhost:8080")
		log.Println("[api] GraphQL Playground: http://localhost:8080/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[api] server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[api] shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[api] shutdown error: %v", err)
	}
}
