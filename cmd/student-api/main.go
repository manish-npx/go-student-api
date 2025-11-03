package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/manish-npx/go-student-api/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Load database (TODO)

	// Setup routes
	route := http.NewServeMux()
	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Welcome to student api"))
	})

	// Setup server
	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: route,
	}

	slog.Info("Server started", slog.String("address", cfg.HttpServer.Addr))

	// Channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	<-done // Block until shutdown signal

	slog.Info("📴 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("❌ Failed to gracefully shutdown server", slog.String("error", err.Error()))
	} else {
		slog.Info("✅ Server shutdown successfully")
	}
}
