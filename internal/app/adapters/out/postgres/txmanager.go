package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type txKey struct{}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(db *Datastore) *TxManager {
	return &TxManager{
		pool: db.pool,
	}
}

func (t *TxManager) ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	err = fn(txCtx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("tx error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func ExtractQuerier(ctx context.Context, pool *pgxpool.Pool) Querier {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if ok {
		return tx
	}
	return pool
}
