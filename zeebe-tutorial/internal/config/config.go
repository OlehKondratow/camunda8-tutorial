package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultAuthorizationServerURL = "https://login.cloud.camunda.io/oauth/token/"
	DefaultWorkerNamePrefix         = "zeebe-tutorial-worker"
)

// Config holds Zeebe client and worker settings.
// Use LocalPlaintext=true (or empty OAuth credentials) for docker-compose Zeebe without TLS/auth.
type Config struct {
	ZeebeAddress           string
	ZeebeClientID          string
	ZeebeClientSecret      string
	AuthorizationServerURL string
	LocalPlaintext         bool
	WorkerName             string
}

// Load reads configuration from the environment. Loads .env when present.
func Load() (*Config, error) {
	_ = godotenv.Load()

	localEnv := strings.TrimSpace(os.Getenv("ZEEBE_LOCAL_PLAINTEXT"))
	if localEnv == "" {
		localEnv = strings.TrimSpace(os.Getenv("ZEEBE_INSECURE_PLAINTEXT"))
	}

	cfg := &Config{
		ZeebeAddress:           strings.TrimSpace(os.Getenv("ZEEBE_ADDRESS")),
		ZeebeClientID:          strings.TrimSpace(os.Getenv("ZEEBE_CLIENT_ID")),
		ZeebeClientSecret:      strings.TrimSpace(os.Getenv("ZEEBE_CLIENT_SECRET")),
		AuthorizationServerURL: strings.TrimSpace(os.Getenv("ZEEBE_AUTHORIZATION_SERVER_URL")),
		WorkerName:             strings.TrimSpace(os.Getenv("WORKER_NAME")),
	}

	if cfg.WorkerName == "" {
		cfg.WorkerName = DefaultWorkerNamePrefix
	}

	if cfg.AuthorizationServerURL == "" {
		cfg.AuthorizationServerURL = DefaultAuthorizationServerURL
		_ = os.Setenv("ZEEBE_AUTHORIZATION_SERVER_URL", cfg.AuthorizationServerURL)
	}

	var errs []error
	if cfg.ZeebeAddress == "" {
		errs = append(errs, errors.New("ZEEBE_ADDRESS is required"))
	}

	oauthMode := cfg.ZeebeClientID != "" || cfg.ZeebeClientSecret != ""
	switch {
	case oauthMode:
		if cfg.ZeebeClientID == "" {
			errs = append(errs, errors.New("ZEEBE_CLIENT_ID is required when using OAuth"))
		}
		if cfg.ZeebeClientSecret == "" {
			errs = append(errs, errors.New("ZEEBE_CLIENT_SECRET is required when using OAuth"))
		}
		if localEnv != "" && strings.EqualFold(localEnv, "true") {
			errs = append(errs, errors.New("ZEEBE_LOCAL_PLAINTEXT must be false (or unset) when using OAuth/TLS"))
		}
		cfg.LocalPlaintext = false
	case localEnv != "":
		cfg.LocalPlaintext = strings.EqualFold(localEnv, "true")
	default:
		cfg.LocalPlaintext = true
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("config: %w", errors.Join(errs...))
	}

	return cfg, nil
}
