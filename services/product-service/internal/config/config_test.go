package config

import (
	"os"
	"testing"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("PRODUCT_SERVICE_HTTP_PORT", "")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.App.Name != "product-service" {
		t.Errorf(
			"expected product-service, got %s",
			cfg.App.Name,
		)
	}

	if cfg.App.Environment != "development" {
		t.Errorf(
			"expected development, got %s",
			cfg.App.Environment,
		)
	}

	if cfg.HTTP.Port != 8080 {
		t.Errorf(
			"expected port 8080, got %d",
			cfg.HTTP.Port,
		)
	}
}

func TestLoadReadsEnvironmentVariables(t *testing.T) {
	t.Setenv("APP_NAME", "gophercart-product")
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PRODUCT_SERVICE_HTTP_PORT", "9090")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.App.Name != "gophercart-product" {
		t.Errorf("unexpected app name: %s", cfg.App.Name)
	}

	if cfg.App.Environment != "production" {
		t.Errorf("unexpected environment: %s", cfg.App.Environment)
	}

	if cfg.HTTP.Port != 9090 {
		t.Errorf("unexpected port: %d", cfg.HTTP.Port)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("PRODUCT_SERVICE_HTTP_PORT", "abc")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
