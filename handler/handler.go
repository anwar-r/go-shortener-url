package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anwar-r/go-shortener-url/redis"

	"github.com/gorilla/mux"
	"github.com/shkshariq/go-util/log"
	"github.com/speps/go-hashids"
)

// ShortenURLRequest represents the request for shortening a URL
type ShortenURLRequest struct {
	URL string `json:"url" example:"https://example.com"`
}

// ShortenURLResponse represents the response after shortening a URL
type ShortenURLResponse struct {
	Shortened string `json:"shortened" example:"http://localhost:8080/abc123"`
}

// ShortenURL godoc
// @Summary Shorten a URL
// @Description Shortens a given URL and returns a shortened version
// @Accept json
// @Produce json
// @Param request body ShortenURLRequest true "URL to shorten"
// @Success 200 {object} ShortenURLResponse
// @Failure 400 {object} map[string]string
// @Router /shorten [post]
func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenURLRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		log.ErrorContext(r.Context(), "Invalid request", err)
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Create a unique short ID for the URL
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	shortID, _ := h.Encode([]int{int(time.Now().Unix())})

	// Store in Redis
	err = redis.Client.Set(context.Background(), shortID, req.URL, 24*time.Hour).Err()
	if err != nil {
		log.ErrorContext(r.Context(), "Failed to store URL", err)
		http.Error(w, `{"error": "Failed to store URL"}`, http.StatusInternalServerError)
		return
	}

	response := ShortenURLResponse{
		Shortened: "http://localhost:8888/" + shortID,
	}
	json.NewEncoder(w).Encode(response)
}

// RedirectURL godoc
// @Summary Redirect to original URL
// @Description Redirects the user to the original URL associated with the shortened ID
// @Param shortID path string true "Shortened ID"
// @Success 302 "Redirects to the original URL"
// @Failure 404 {object} map[string]string
// @Router /{shortID} [get]
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortID"]

	// Retrieve from Redis
	originalURL, err := redis.Client.Get(context.Background(), shortID).Result()
	if err != nil {
		log.ErrorContext(r.Context(), "URL not found", err)
		http.Error(w, `{"error": "URL not found"}`, http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}
