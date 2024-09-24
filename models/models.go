package models

// ShortenURLRequest is used for the /shorten API
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse is the response for the /shorten API
type ShortenURLResponse struct {
	Shortened string `json:"shortened"`
}
