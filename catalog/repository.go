package catalog

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListsProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListsProductsWithIDs(ctx context.Context)
	SearchPrducts(ctx context.Context)
}

type elasticRepository struct {
	client *elastic.Client
}

type ProductDocument struct {
	Name
	Description
	Price
}
