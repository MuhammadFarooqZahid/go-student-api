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
	"github.com/MuhammadFarooqZahid/go-student-api/internal/http/handler/student"
	"github.com/MuhammadFarooqZahid/go-student-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env))

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentsList(storage))

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
