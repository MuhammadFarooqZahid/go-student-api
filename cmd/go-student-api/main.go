package main

import (
	"log"
	"net/http"

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

	log.Printf("Server is running at %s\n", cfg.Address)

	err := http.ListenAndServe(cfg.Address, router)

	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

}
