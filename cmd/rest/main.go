package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frkntplglu/insider/internal/container"
	"github.com/frkntplglu/insider/pkg/logger"
)

func main() {
	container := container.NewContainer()

	go func() {
		if err := container.Start(); err != nil {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := container.Stop(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
