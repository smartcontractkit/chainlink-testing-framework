package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/linkingservice"
)

func main() {
	server := linkingservice.NewServer(linkingservice.Config{
		GRPCPort:  envInt("LINKING_SERVICE_GRPC_PORT", linkingservice.DefaultGRPCPort),
		AdminPort: envInt("LINKING_SERVICE_ADMIN_PORT", linkingservice.DefaultAdminPort),
	})

	if err := server.Start(context.Background()); err != nil {
		panic(err)
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
