package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config はPostgreSQLの接続パラメータを保持する。
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DefaultConfig はローカル開発用の設定を返す。
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "delivery",
		Password: "delivery",
		DBName:   "food_delivery",
		SSLMode:  "disable",
	}
}

// DSN は接続文字列を返す。
func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// NewPool はConfigからコネクションプールを作成する。
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
