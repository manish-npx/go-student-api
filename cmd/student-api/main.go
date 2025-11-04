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
	"github.com/manish-npx/go-student-api/internal/storage/factory"
)

func main() {
	// 🧩 Load config
	cfg := config.MustLoad()

	// Load database (COMPLETED)

	/// 🧩 Choose database based on config
	//var storage storage.Storage
	var err error

	storage, err := factory.NewStorage(*cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// 🧩 Setup routes
	route := http.NewServeMux()
	route.HandleFunc("POST /api/student", student.New(storage))
	route.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	route.HandleFunc("GET /api/students", student.GetList(storage))
	route.HandleFunc("PUT /api/student/{id}", student.UpdateById(storage))

	// 🧩 Setup server
	server := &http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: route,
	}

	slog.Info("Server started", slog.String("address", cfg.HttpServer.Addr))

	slog.Info("💾 Database initialized",
		slog.String("driver", cfg.DBType),
	)

	// Channel for graceful shutdown
	// 🧩 Graceful shutdown
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
