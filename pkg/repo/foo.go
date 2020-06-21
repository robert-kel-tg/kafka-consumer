package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/robertke/kafka-consumer/pkg/domain"
)

type (
	FooRepo struct {
		db *sqlx.DB
	}
)

func NewFoo(db *sqlx.DB) *FooRepo {
	return &FooRepo{db}
}

func (r FooRepo) Insert(ctx context.Context, tx *sqlx.Tx, msg domain.Foo) (*domain.Foo, error) {

	query := `INSERT INTO messages 
	(	name,
	) 
		VALUES ($1) 
		RETURNING *`

	var (
		foo domain.Foo
	)

	if err := tx.GetContext(ctx, &foo, query, msg); err != nil {
		return &domain.Foo{}, err
	}

	return &foo, nil
}

func (r FooRepo) Update(ctx context.Context, tx *sqlx.Tx) (*domain.Foo, error) {
	return nil, nil
}

func (r FooRepo) GetOne(ctx context.Context, id string) (*domain.Foo, error) {
	return nil, nil
}
