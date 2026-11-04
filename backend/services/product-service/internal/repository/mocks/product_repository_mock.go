package mocks

import (
	"context"
	"sort"
	"sync"

	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/domain/product"
	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/repository"
)

type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]*product.Product

	CreateFunc func(context.Context, *product.Product) error
	ListFunc   func(context.Context, int, int) ([]*product.Product, error)
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]*product.Product),
	}
}

func (m *ProductRepository) Create(
	ctx context.Context,
	p *product.Product,
) error {
	if m.CreateFunc != nil {
		if err := m.CreateFunc(ctx, p); err != nil {
			return err
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	copied := *p
	m.products[p.ID] = &copied

	return nil
}

func (m *ProductRepository) GetByID(
	ctx context.Context,
	id string,
) (*product.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.products[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	copied := *p
	return &copied, nil
}

func (m *ProductRepository) List(
	ctx context.Context,
	limit int,
	offset int,
) ([]*product.Product, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, limit, offset)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.products))
	for id := range m.products {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	if offset >= len(ids) {
		return []*product.Product{}, nil
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}

	result := make([]*product.Product, 0, end-offset)

	for _, id := range ids[offset:end] {
		copied := *m.products[id]
		result = append(result, &copied)
	}

	return result, nil
}

func (m *ProductRepository) Update(
	ctx context.Context,
	p *product.Product,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[p.ID]; !exists {
		return repository.ErrNotFound
	}

	copied := *p
	m.products[p.ID] = &copied

	return nil
}

func (m *ProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[id]; !exists {
		return repository.ErrNotFound
	}

	delete(m.products, id)
	return nil
}
