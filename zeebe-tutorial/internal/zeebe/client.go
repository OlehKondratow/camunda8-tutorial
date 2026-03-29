package zeebe

import (
	"fmt"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"zeebe-tutorial/internal/config"
)

// NewClient builds a Zeebe gRPC client from config (OAuth for Camunda SaaS or plaintext for local broker).
func NewClient(cfg *config.Config) (*zbc.Client, error) {
	cc := &zbc.ClientConfig{
		GatewayAddress: cfg.ZeebeAddress,
		UsePlaintext:   cfg.LocalPlaintext,
	}
	if cfg.LocalPlaintext && (cfg.ZeebeClientID != "" || cfg.ZeebeClientSecret != "") {
		return nil, fmt.Errorf("zeebe client: plaintext mode cannot be used with OAuth credentials set")
	}
	client, err := zbc.NewClient(cc)
	if err != nil {
		return nil, fmt.Errorf("zeebe new client: %w", err)
	}
	return client, nil
}
