package ports

import "context"

type TxManager interface {
	ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error
}
