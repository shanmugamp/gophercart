package product

import "testing"

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		description string
		priceMinor  int64
		currency    string
		stock       int32
		wantErr     bool
	}{
		{
			name:        "valid product",
			productName: "iPhone 17",
			description: "Apple smartphone",
			priceMinor:  99900,
			currency:    "INR",
			stock:       10,
		},
		{
			name:        "empty product name",
			description: "Apple smartphone",
			priceMinor:  99900,
			currency:    "INR",
			stock:       10,
			wantErr:     true,
		},
		{
			name:        "negative price",
			productName: "iPhone 17",
			description: "Apple smartphone",
			priceMinor:  -1,
			currency:    "INR",
			stock:       10,
			wantErr:     true,
		},
		{
			name:        "negative stock",
			productName: "iPhone 17",
			description: "Apple smartphone",
			priceMinor:  99900,
			currency:    "INR",
			stock:       -1,
			wantErr:     true,
		},
		{
			name:        "empty currency",
			productName: "iPhone 17",
			description: "Apple smartphone",
			priceMinor:  99900,
			stock:       10,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.description, tt.priceMinor, tt.currency, tt.stock)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if product == nil {
				t.Fatal("expected product but got nil")
			}
			if product.Name != tt.productName {
				t.Errorf("expected product name %q but got %q", tt.productName, product.Name)
			}
			if product.Description != tt.description {
				t.Errorf("expected product description %q but got %q", tt.description, product.Description)
			}
			if product.PriceMinor != tt.priceMinor {
				t.Errorf("expected product price %d but got %d", tt.priceMinor, product.PriceMinor)
			}
			if product.Currency != tt.currency {
				t.Errorf("expected product currency %q but got %q", tt.currency, product.Currency)
			}
			if product.Stock != tt.stock {
				t.Errorf("expected product stock %d but got %d", tt.stock, product.Stock)
			}
			if product.ID == "" {
				t.Error("expected product ID to be set but got empty string")
			}
			if !product.Active {
				t.Error("New product should be active")
			}
		})
	}
}

func TestProductDecreaseStock(t *testing.T) {
	product, err := NewProduct("iPhone 17", "Apple smartphone", 99900, "INR", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = product.DecreaseStock(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.Stock != 5 {
		t.Errorf("expected product stock 5 but got %d", product.Stock)
	}
}

func TestProductDecreaseStockInsufficientStock(t *testing.T) {
	product, err := NewProduct(
		"Laptop",
		"Business laptop",
		7500000,
		"INR",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = product.DecreaseStock(11)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if product.Stock != 10 {
		t.Fatalf(
			"stock should remain unchanged, got %d",
			product.Stock,
		)
	}
}
func TestProductIncreaseStock(t *testing.T) {
	product, err := NewProduct(
		"Laptop",
		"Business laptop",
		7500000,
		"INR",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = product.IncreaseStock(5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.Stock != 15 {
		t.Fatalf("expected stock 15, got %d", product.Stock)
	}
}
func TestProductDeactivate(t *testing.T) {
	product, err := NewProduct(
		"Laptop",
		"Business laptop",
		7500000,
		"INR",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	product.Deactivate()

	if product.Active {
		t.Fatal("expected product to be inactive")
	}
}
func TestProductActivate(t *testing.T) {
	product, err := NewProduct(
		"Laptop",
		"Business laptop",
		7500000,
		"INR",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	product.Deactivate()
	product.Activate()

	if !product.Active {
		t.Fatal("expected product to be active")
	}
}
