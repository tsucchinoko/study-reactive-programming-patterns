package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql/generated"
	"github.com/gorilla/websocket"
)

// NewHandler creates the GraphQL HTTP handler with WebSocket support for subscriptions.
func NewHandler(resolver *Resolver) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	})

	return srv
}

// NewPlaygroundHandler creates the GraphQL Playground UI handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("Food Delivery Tracker", endpoint)
}
