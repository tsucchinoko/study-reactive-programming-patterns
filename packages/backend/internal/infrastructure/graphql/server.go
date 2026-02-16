package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daichitsuchiya/food-delivery-tracker/internal/infrastructure/graphql/generated"
)

// NewHandler creates the GraphQL HTTP handler.
func NewHandler(resolver *Resolver) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	return srv
}

// NewPlaygroundHandler creates the GraphQL Playground UI handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("Food Delivery Tracker", endpoint)
}
