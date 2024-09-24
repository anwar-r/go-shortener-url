package main

import (
	"log"
	"net/http"

	_ "go-shortener-url/docs"
	"go-shortener-url/handler"
	"go-shortener-url/redis"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/joho/godotenv"
)

// @title URL Shortener API
// @version 1.0
// @description This is a simple URL shortener service.
// @host localhost:8888
// @BasePath /

func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Redis connection
	redis.Initialize()

	// Initialize router
	r := mux.NewRouter()

	// API Routes
	r.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortID}", handler.RedirectURL).Methods("GET")

	// Swagger docs route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start server
	log.Println("Server started on :8888")
	log.Fatal(http.ListenAndServe(":8888", r))
}
