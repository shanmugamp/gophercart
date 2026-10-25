package product

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNameRequired      = errors.New("product name is required")
	ErrPriceInvalid      = errors.New("product price cannot be negative")
	ErrCurrencyRequired  = errors.New("product currency is required")
	ErrStockInvalid      = errors.New("product stock cannot be negative")
	ErrInsufficientStock = errors.New("insufficient product stock")
)

type Product struct {
	ID          string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Stock       int32
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(name string, description string, priceMinor int64, currency string, stock int32) (*Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%s", ErrNameRequired)
	}

	if priceMinor < 0 {
		return nil, fmt.Errorf("%s", ErrPriceInvalid)
	}
	currency = strings.TrimSpace(currency)
	if currency == "" {
		return nil, fmt.Errorf("%s", ErrCurrencyRequired)
	}
	if stock < 0 {
		return nil, fmt.Errorf("%s", ErrStockInvalid)
	}
	now := time.Now().UTC()
	return &Product{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		PriceMinor:  priceMinor,
		Currency:    currency,
		Stock:       stock,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
func (p *Product) IncreaseStock(quantity int32) error {
	if quantity <= 0 {
		return ErrStockInvalid
	}

	p.Stock += quantity
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func (p *Product) DecreaseStock(quantity int32) error {
	if quantity <= 0 {
		return ErrStockInvalid
	}
	if p.Stock < quantity {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now().UTC()
	return nil
}
func (p *Product) Activate() {
	p.Active = true
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) Deactivate() {
	p.Active = false
	p.UpdatedAt = time.Now().UTC()
}
