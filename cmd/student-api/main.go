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

// load config
// load database
// setup server
// setup route
func main() {
	//load config
	cfg := config.MustLoad()

	//load database

	//setup route
	route := http.NewServeMux()
	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		//fmt.Fprintf(w, "Welcome to the home page!!  ")
		w.Write([]byte("Welcome to student api"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: route,
	}

	log.Printf("Server started on Addr %s", cfg.HttpServer.Addr)

	//channel signal
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGILL, syscall.SIGTERM)

	go func() {

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting down server ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")

}
