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
	"github.com/manish-npx/go-student-api/internal/http/handlers/student"
	"github.com/manish-npx/go-student-api/internal/storage/postgres"
	//"github.com/manish-npx/go-student-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Load database (COMPLETED)

	//sqlite loading
	//storage, err := sqlite.New(*cfg)

	//pgsql
	storage, err := postgres.New(*cfg)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	if err != nil {
		log.Fatal("Error! database connection issue", err)
	}
	slog.Info("Storage init", slog.String("env", cfg.Env))

	slog.Info("Storage")

	// Setup routes
	route := http.NewServeMux()
	route.HandleFunc("POST /api/student", student.New(storage))
	route.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	route.HandleFunc("GET /api/students", student.GetList(storage))

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
