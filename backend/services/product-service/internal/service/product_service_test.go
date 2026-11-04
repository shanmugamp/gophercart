package service

import (
	"context"
	"errors"
	"testing"

	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/domain/product"
	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/repository"
	"github.com/shanmugamp/gophercart/backend/services/product-service/internal/repository/mocks"
)

func TestCreateProduct(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	ctx := context.Background()

	created, err := svc.CreateProduct(
		ctx,
		"Laptop",
		"Developer laptop",
		8500000,
		"INR",
		10,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created == nil {
		t.Fatal("expected product, got nil")
	}

	if created.Name != "Laptop" {
		t.Errorf("expected Laptop, got %s", created.Name)
	}

	if created.PriceMinor != 8500000 {
		t.Errorf("expected price 8500000, got %d", created.PriceMinor)
	}

	if created.Stock != 10 {
		t.Errorf("expected stock 10, got %d", created.Stock)
	}

	if created.ID == "" {
		t.Error("expected product ID to be generated")
	}
}

func TestCreateProductInvalidName(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	_, err := svc.CreateProduct(
		context.Background(),
		"   ",
		"Description",
		10000,
		"INR",
		10,
	)

	if !errors.Is(err, product.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestCreateProductRepositoryError(t *testing.T) {
	repo := mocks.NewProductRepository()
	expectedErr := errors.New("database unavailable")

	repo.CreateFunc = func(
		ctx context.Context,
		p *product.Product,
	) error {
		return expectedErr
	}

	svc := NewProductService(repo)

	_, err := svc.CreateProduct(
		context.Background(),
		"Laptop",
		"Developer laptop",
		8500000,
		"INR",
		10,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestGetProduct(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	created, err := svc.CreateProduct(
		context.Background(),
		"Laptop",
		"Developer laptop",
		8500000,
		"INR",
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.GetProduct(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestGetProductNotFound(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	_, err := svc.GetProduct(context.Background(), "missing-id")

	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListProductsPagination(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	_, err := svc.CreateProduct(
		context.Background(), "Laptop", "Laptop", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.CreateProduct(
		context.Background(), "Mouse", "Mouse", 2000, "INR", 10,
	)
	if err != nil {
		t.Fatal(err)
	}

	products, err := svc.ListProducts(context.Background(), 1, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

func TestListProductsDefaultLimit(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	var receivedLimit int

	repo.ListFunc = func(
		ctx context.Context,
		limit, offset int,
	) ([]*product.Product, error) {
		receivedLimit = limit
		return nil, nil
	}

	_, err := svc.ListProducts(context.Background(), 0, 0)

	if err != nil {
		t.Fatal(err)
	}

	if receivedLimit != 20 {
		t.Errorf("expected default limit 20, got %d", receivedLimit)
	}
}

func TestListProductsMaximumLimit(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	var receivedLimit int

	repo.ListFunc = func(
		ctx context.Context,
		limit, offset int,
	) ([]*product.Product, error) {
		receivedLimit = limit
		return nil, nil
	}

	_, err := svc.ListProducts(context.Background(), 500, 0)

	if err != nil {
		t.Fatal(err)
	}

	if receivedLimit != 100 {
		t.Errorf("expected maximum limit 100, got %d", receivedLimit)
	}
}

func TestUpdateProduct(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(
		context.Background(), "Laptop", "Old description", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	p.Description = "Updated description"

	err = svc.UpdateProduct(context.Background(), p)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := svc.GetProduct(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Description != "Updated description" {
		t.Errorf("product description was not updated")
	}
}

func TestDeleteProduct(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(
		context.Background(), "Laptop", "Laptop", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.DeleteProduct(context.Background(), p.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.GetProduct(context.Background(), p.ID)

	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestIncreaseStock(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(
		context.Background(), "Laptop", "Laptop", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.IncreaseStock(context.Background(), p.ID, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := svc.GetProduct(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Stock != 8 {
		t.Errorf("expected stock 8, got %d", updated.Stock)
	}
}

func TestDecreaseStock(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(
		context.Background(), "Laptop", "Laptop", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.DecreaseStock(context.Background(), p.ID, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := svc.GetProduct(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Stock != 3 {
		t.Errorf("expected stock 3, got %d", updated.Stock)
	}
}

func TestDecreaseStockInsufficient(t *testing.T) {
	repo := mocks.NewProductRepository()
	svc := NewProductService(repo)

	p, err := svc.CreateProduct(
		context.Background(), "Laptop", "Laptop", 10000, "INR", 5,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.DecreaseStock(context.Background(), p.ID, 10)

	if !errors.Is(err, product.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}
