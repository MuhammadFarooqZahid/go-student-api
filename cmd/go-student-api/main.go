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

	"github.com/MuhammadFarooqZahid/go-student-api/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup
	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student's api"))
	})

	// Setup server
	slog.Info("Server is running at", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
