package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"multimodel-db-engine/internal/config"
	"multimodel-db-engine/internal/server"
	"multimodel-db-engine/internal/database"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database engine
	dbEngine := database.NewMultiModelDatabase(cfg)

	// Create HTTP router
	router := mux.NewRouter()

	// Setup API routes
	server.SetupRoutes(router, dbEngine)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(router)

	log.Printf("Starting Multi-Model Database Engine on port %s", cfg.Port)
	
	// Start the server
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}