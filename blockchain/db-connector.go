package blockchain

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, transaction Transaction) (*Transaction, error)
	All(ctx context.Context) ([]Transaction, error)
	GetByName(ctx context.Context, name string) (*Transaction, error)
	Update(ctx context.Context, id int64, updated Transaction) (*Transaction, error)
	Delete(ctx context.Context, id int64) error
}
