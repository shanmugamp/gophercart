package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
}

type AppConfig struct {
	Name        string
	Environment string
	LogLevel    string
}

type HTTPConfig struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	port, err := getInt("PRODUCT_SERVICE_HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}

	return Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "product-service"),
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		HTTP: HTTPConfig{
			Host:              getEnv("HTTP_HOST", "0.0.0.0"),
			Port:              port,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ShutdownTimeout:   10 * time.Second,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid integer value for %s: %w",
			key,
			err,
		)
	}

	return result, nil
}
