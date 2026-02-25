package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisConfig は Redis 接続設定。
type RedisConfig struct {
	Addr string
}

// DefaultConfig は環境変数 REDIS_ADDR またはデフォルト localhost:6379 を返す。
func DefaultConfig() RedisConfig {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return RedisConfig{Addr: addr}
}

// NewClient は Redis クライアントを生成する。
func NewClient(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
}
