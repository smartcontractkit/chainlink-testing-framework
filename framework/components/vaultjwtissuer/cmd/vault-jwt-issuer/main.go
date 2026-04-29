package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	vaultjwtissuer "github.com/smartcontractkit/chainlink-testing-framework/framework/components/vaultjwtissuer"
)

func main() {
	httpPort := envInt("VAULT_JWT_ISSUER_HTTP_PORT", vaultjwtissuer.DefaultHTTPPort)

	server, err := vaultjwtissuer.NewServer(vaultjwtissuer.Config{
		HTTPPort: httpPort,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create vault JWT issuer: %w", err))
	}

	if err := server.Start(context.Background()); err != nil {
		panic(fmt.Errorf("failed to start vault JWT issuer: %w", err))
	}
	defer func() {
		_ = server.Stop()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

func envInt(name string, defaultValue int) int {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
