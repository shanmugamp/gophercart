package repository

import (
	"context"

	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/domain/product"
)

type ProductRepository interface {
	Create(ctx context.Context, p *product.Product) error

	GetByID(ctx context.Context, id string) (*product.Product, error)

	List(ctx context.Context, limit, offset int) ([]*product.Product, error)

	Update(ctx context.Context, p *product.Product) error

	Delete(ctx context.Context, id string) error
}
