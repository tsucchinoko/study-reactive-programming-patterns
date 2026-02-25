package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTransaction はデータベーストランザクション内でfnを実行する。
// fnがエラーを返した場合、トランザクションはロールバックされる。
// 関数型ラッパー — 呼び出し側が純粋なロジックを渡し、この関数が副作用を処理する。
func WithTransaction[T any](ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) (T, error)) (T, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("begin transaction: %w", err)
	}

	result, err := fn(tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		var zero T
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		var zero T
		return zero, fmt.Errorf("commit transaction: %w", err)
	}

	return result, nil
}
