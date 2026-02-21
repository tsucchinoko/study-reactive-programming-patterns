package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/tsucchinoko/food-delivery-tracker/internal/infrastructure/graphql/generated"
)

// NewHandler はWebSocketサブスクリプション対応のGraphQL HTTPハンドラを作成する。
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
				return true // 開発環境では全オリジンを許可
			},
		},
	})

	return srv
}

// NewPlaygroundHandler はGraphQL Playground UIハンドラを作成する。
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("Food Delivery Tracker", endpoint)
}
