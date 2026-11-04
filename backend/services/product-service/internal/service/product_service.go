package service

import (
	"context"
	"errors"

	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/domain/product"
	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/repository"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrInvalidProduct = errors.New("invalid product")
)

type ProductService struct {
	repository repository.ProductRepository
}

func NewProductService(
	repo repository.ProductRepository,
) *ProductService {
	return &ProductService{
		repository: repo,
	}
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(
	ctx context.Context,
	name string,
	description string,
	priceMinor int64,
	currency string,
	stock int32,
) (*product.Product, error) {

	p, err := product.NewProduct(
		name,
		description,
		priceMinor,
		currency,
		stock,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// GetProduct retrieves a product by ID.
func (s *ProductService) GetProduct(
	ctx context.Context,
	id string,
) (*product.Product, error) {
	return s.repository.GetByID(ctx, id)
}

// ListProducts retrieves products with pagination.
func (s *ProductService) ListProducts(
	ctx context.Context,
	limit int,
	offset int,
) ([]*product.Product, error) {

	if limit <= 0 {
		limit = defaultPageSize
	}

	if limit > maxPageSize {
		limit = maxPageSize
	}

	if offset < 0 {
		offset = 0
	}

	return s.repository.List(ctx, limit, offset)
}

// UpdateProduct updates an existing product.
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	p *product.Product,
) error {

	if p == nil {
		return ErrInvalidProduct
	}

	return s.repository.Update(ctx, p)
}

// DeleteProduct deletes a product by ID.
func (s *ProductService) DeleteProduct(
	ctx context.Context,
	id string,
) error {
	return s.repository.Delete(ctx, id)
}

// IncreaseStock increases the stock of a product.
func (s *ProductService) IncreaseStock(
	ctx context.Context,
	id string,
	quantity int32,
) error {

	p, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.IncreaseStock(quantity); err != nil {
		return err
	}

	return s.repository.Update(ctx, p)
}

// DecreaseStock decreases the stock of a product.
func (s *ProductService) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int32,
) error {

	p, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.DecreaseStock(quantity); err != nil {
		return err
	}

	return s.repository.Update(ctx, p)
}
